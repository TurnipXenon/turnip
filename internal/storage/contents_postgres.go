package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"

	"github.com/TurnipXenon/turnip/internal/storage/migration"
	"github.com/TurnipXenon/turnip/internal/util"
)

type contentsPostgresImpl struct {
	db        *PostgresDb
	tableName string
	tags      Tags
}

func (c *contentsPostgresImpl) GetTableName() string {
	return c.tableName
}

func (c *contentsPostgresImpl) GetMigrationSequence() []migration.Migration {
	return []migration.Migration{
		migration.NewGenericMigration(migration.MigrateContent0001),
	}
}

func (c *contentsPostgresImpl) CreateContent(ctx context.Context, request *turnip.ContentRequestResponse, user *turnip.User) (*turnip.Content, error) {
	// todo: require some fields!

	// create uuid
	// very unlikely to collide, right?
	content := request.Item
	content.PrimaryId = uuid.New().String()
	content.CreatedAt = timestamppb.Now()
	content.AuthorId = user.PrimaryId

	meta, err := json.Marshal(content.Meta)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	_, err = c.db.Pool.Exec(ctx, `INSERT INTO "Content"
    	(primary_id, created_at, title, description, content, tag_list, meta, author_id, slug) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		content.PrimaryId, content.CreatedAt.AsTime().Format(time.RFC3339), // primary_id: $1, created_at: $2
		content.Title, content.Description, content.Content, // title: $3, description: $4, content: $5
		content.TagList, meta, // tag_list: $6, meta: $7
		content.AuthorId, content.Slug, // author_id: $8, slug: $9 todo(turnip): missing media!
	)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	err = c.tags.UpdateTags(ctx, content)
	if err != nil {
		util.LogDetailedError(err)
	}

	return content, nil
}

func pgxUuidToGoogleUuid(initial pgtype.UUID) (*uuid.UUID, error) {
	final, err := uuid.FromBytes(initial.Bytes[:])
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}
	return &final, nil
}

func pgxUuidToStringUuid(initial pgtype.UUID) (string, error) {
	final, err := pgxUuidToGoogleUuid(initial)
	if err != nil {
		util.LogDetailedError(err)
		return "", err
	}
	return final.String(), nil
}

func stringToStringPtr(value string) *string {
	tmp := value
	return &tmp
}

// GetContentById returns nil content also with nil error!
// todo: document behavior
func (c *contentsPostgresImpl) GetContentById(ctx context.Context, idQuery string) (*turnip.Content, error) {
	if idQuery == "" {
		return nil, fmt.Errorf("id is empty")
	}

	row := c.db.Pool.QueryRow(ctx, `SELECT t.*
               FROM "Content" t
               WHERE primary_id = $1
               LIMIT 1`, idQuery)

	var content turnip.Content
	var primaryId pgtype.UUID
	var authorId pgtype.UUID
	var createdAt pgtype.Timestamp
	var title, description, contentString, accessDetails, meta, slugString *string
	// todo: turn to CollectRow
	err := row.Scan(&primaryId, &createdAt, &title, &description, &contentString,
		&content.TagList, &accessDetails, &meta, &authorId, &slugString)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	content.Title = derefString(title)
	content.Description = derefString(description)
	content.Content = derefString(contentString)
	content.Slug = slugString
	content.PrimaryId, err = pgxUuidToStringUuid(primaryId)
	content.AuthorId, err = pgxUuidToStringUuid(authorId)
	content.CreatedAt = timestamppb.New(createdAt.Time)
	if meta == nil {
		meta = stringToStringPtr("")
	}
	json.Unmarshal([]byte(*meta), &content.Meta)
	//content.Meta = map[string]string{}

	content.TagList, err = c.tags.GetTagsByContent(ctx, &content)
	if err != nil {
		util.LogDetailedError(err)
	}

	return &content, nil
}

func (c *contentsPostgresImpl) GetContentBySlug(ctx context.Context, slug string) (*turnip.Content, error) {
	if slug == "" {
		return nil, fmt.Errorf("slug is empty")
	}

	row := c.db.Pool.QueryRow(ctx, `SELECT t.*
               FROM "Content" t
               WHERE slug = $1
               LIMIT 1`, slug)

	var content turnip.Content
	var primaryId pgtype.UUID
	var authorId pgtype.UUID
	var createdAt pgtype.Timestamp
	var title, description, contentString, accessDetails, meta, slugString *string
	// todo: turn to CollectRow
	err := row.Scan(&primaryId, &createdAt, &title, &description, &contentString,
		&content.TagList, &accessDetails, &meta, &authorId, &slugString)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	content.Title = derefString(title)
	content.Description = derefString(description)
	content.Content = derefString(contentString)
	content.Slug = slugString
	content.PrimaryId, err = pgxUuidToStringUuid(primaryId)
	content.AuthorId, err = pgxUuidToStringUuid(authorId)
	content.CreatedAt = timestamppb.New(createdAt.Time)
	// todo: parse from string accessDetails and meta
	if meta == nil {
		meta = stringToStringPtr("")
	}
	json.Unmarshal([]byte(*meta), &content.Meta)
	//content.Meta = map[string]string{}

	content.TagList, err = c.tags.GetTagsByContent(ctx, &content)
	if err != nil {
		util.LogDetailedError(err)
	}

	return &content, nil
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func (c *contentsPostgresImpl) rowsToContentList(ctx context.Context, rows pgx.Rows) ([]*turnip.Content, error) {
	contentList := []*turnip.Content{}
	var primaryId, authorId pgtype.UUID
	var createdAt pgtype.Timestamp
	var title, description, content, accessDetails, meta, slugString *string
	var tagList []string
	_, err := pgx.ForEachRow(rows, []any{&primaryId, &createdAt, &title, &description, &content,
		&tagList, &accessDetails, &meta, &authorId, &slugString}, func() error {
		// todo: check if accessible, otherwise add to list
		newContent := &turnip.Content{
			Title:       derefString(title),
			Description: derefString(description),
			Content:     derefString(content),
			Slug:        slugString,
			TagList:     tagList,
			Meta:        nil, // todo parse
			// todo media field
		}

		var err error
		newContent.PrimaryId, err = pgxUuidToStringUuid(primaryId)
		newContent.AuthorId, err = pgxUuidToStringUuid(authorId)
		if err != nil {
			util.LogDetailedError(err)
		}

		newContent.CreatedAt = timestamppb.New(createdAt.Time)
		newContent.TagList, err = c.tags.GetTagsByContent(ctx, newContent)
		if meta == nil {
			meta = stringToStringPtr("")
		}
		json.Unmarshal([]byte(*meta), &newContent.Meta)
		contentList = append(contentList, newContent)
		return nil
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}
	return contentList, nil
}

func (c *contentsPostgresImpl) GetAllContent(ctx context.Context) ([]*turnip.Content, error) {
	rows, _ := c.db.Pool.Query(ctx, `SELECT * FROM "Content"`)
	return c.rowsToContentList(ctx, rows)
}

func (c *contentsPostgresImpl) GetContentByTagInclusive(ctx context.Context, tag []string) ([]*turnip.Content, error) {
	contentIdList, err := c.tags.GetContentIdsByTagInclusive(ctx, tag)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}
	if len(contentIdList) == 0 {
		return nil, nil
	}

	// todo: might need normalization
	vl := convertToAnyList(contentIdList)
	q := stringListToSanitizedSql(vl, 0)
	rows, _ := c.db.Pool.Query(ctx,
		fmt.Sprintf(`SELECT * FROM "%s" WHERE primary_id IN (%s)`,
			c.tableName, q), vl...)
	return c.rowsToContentList(ctx, rows)
}

func (c *contentsPostgresImpl) UpdateContent(ctx context.Context, newContent *turnip.Content) (*turnip.Content, error) {
	// todo: make setting more dynamic instead of setting everything
	// todo set these attributes to UpdateContent
	accessDetails := &turnip.AccessDetails{}
	meta, err := json.Marshal(newContent.Meta)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	_, err = c.db.Pool.Exec(ctx, `UPDATE public."Content"
		SET title=$1, description=$2, content=$3, tag_list=$4, access_details=$5, meta=$6, slug=$8
		WHERE primary_id = $7`,
		newContent.Title, newContent.Description, newContent.Content, // 1-3
		newContent.TagList, accessDetails, meta, newContent.PrimaryId, // 4-7
		newContent.Slug, // 8
	)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	// todo: put NoRowErr check here!
	err = c.tags.UpdateTags(ctx, newContent)
	if err != nil {
		util.LogDetailedError(err)
	}

	return newContent, nil
}

func (c *contentsPostgresImpl) DeleteContentById(ctx context.Context, primaryId string) (*turnip.Content, error) {
	_, err := c.db.Pool.Exec(ctx, `DELETE FROM "Content" 
       WHERE primary_id = $1`, primaryId)

	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	return nil, nil
}

func (c *contentsPostgresImpl) GetContentByTagStrict(ctx context.Context, tag []string) ([]*turnip.Content, error) {
	contentIdList, err := c.tags.GetContentIdsByTagStrict(ctx, tag) // todo
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}
	if len(contentIdList) == 0 {
		return nil, nil
	}

	// todo: normalize
	vl := convertToAnyList(contentIdList)
	q := stringListToSanitizedSql(vl, 0)
	rows, _ := c.db.Pool.Query(ctx,
		fmt.Sprintf(`SELECT * FROM "%s" WHERE primary_id IN (%s)`,
			c.tableName, q), vl...)
	return c.rowsToContentList(ctx, rows)
}

func NewContentsPostgres(ctx context.Context, d *PostgresDb, tags Tags) Contents {
	// primary: primary id
	// sort: created at
	t := contentsPostgresImpl{
		db:        d,
		tableName: "Content",
		tags:      tags,
	}

	SetupTable(ctx, d, &t)

	return &t
}
