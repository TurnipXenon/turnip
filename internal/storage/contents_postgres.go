package storage

import (
	"context"
	"encoding/json"
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
	// todo: global secondary index
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

	accessDetails, err := json.Marshal(content.AccessDetails)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}
	meta, err := json.Marshal(content.Meta)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	_, err = c.db.Pool.Exec(ctx, `INSERT INTO "Content"
    	(primary_id, created_at, title, description, content, tag_list, access_details, meta, author_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		content.PrimaryId, content.CreatedAt.AsTime().Format(time.RFC3339), // primary_id: $1, created_at: $2
		content.Title, content.Description, content.Content, // title: $3, description: $4, content: $5
		content.TagList, accessDetails, meta, // tag_list: $6, access_details: $7, meta: $8, author_id: $9
		content.AuthorId, // author_id: $9, todo(turnip): missing media!
	)
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
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

func pgxByteToStringUuid(initial []byte) (string, error) {
	final, err := uuid.FromBytes(initial)
	if err != nil {
		util.LogDetailedError(err)
		return "", util.WrapErrorWithDetails(err)
	}
	return final.String(), nil
}

// GetContentById returns nil content also with nil error!
// todo: document behavior
func (c *contentsPostgresImpl) GetContentById(ctx context.Context, idQuery string) (*turnip.Content, error) {
	row := c.db.Pool.QueryRow(ctx, `SELECT t.*
               FROM "Content" t
               WHERE primary_id = $1
               LIMIT 1`, idQuery)

	var content turnip.Content
	var primaryId pgtype.UUID
	var authorId pgtype.UUID
	var createdAt pgtype.Timestamp
	var accessDetails, meta string
	// todo: turn to CollectRow
	err := row.Scan(&primaryId, &createdAt, &content.Title, &content.Description, &content.Content,
		&content.TagList, &accessDetails, &meta, &authorId)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	content.PrimaryId, err = pgxUuidToStringUuid(primaryId)
	content.AuthorId, err = pgxUuidToStringUuid(authorId)
	content.CreatedAt = timestamppb.New(createdAt.Time)
	// todo: parse from string accessDetails and meta
	content.AccessDetails = &turnip.AccessDetails{}
	content.Meta = map[string]string{}

	return &content, nil
}

func (c *contentsPostgresImpl) GetAllContent(ctx context.Context) ([]*turnip.Content, error) {
	contentList := []*turnip.Content{}

	// todo: handle paging
	var primaryId, authorId pgtype.UUID
	var createdAt pgtype.Timestamp
	var title, description, content string
	var tagList []string
	var accessDetails, meta string

	rows, _ := c.db.Pool.Query(ctx, `SELECT * FROM "Content"`)
	_, err := pgx.ForEachRow(rows, []any{&primaryId, &createdAt, &title, &description, &content,
		&tagList, &accessDetails, &meta, &authorId}, func() error {
		// todo: check if access, otherwise add to list
		newContent := &turnip.Content{
			Title:         title,
			Description:   description,
			Content:       content,
			TagList:       tagList,
			AccessDetails: nil, // todo parse
			Meta:          nil, // todo parse
			// todo media field
		}

		var err error
		newContent.PrimaryId, err = pgxUuidToStringUuid(primaryId)
		newContent.AuthorId, err = pgxUuidToStringUuid(authorId)
		if err != nil {
			util.LogDetailedError(err)
		}

		newContent.CreatedAt = timestamppb.New(createdAt.Time)
		contentList = append(contentList, newContent)
		return nil
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	return contentList, nil
}

func (c *contentsPostgresImpl) GetContentByTag(ctx context.Context, tag string) ([]*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c *contentsPostgresImpl) UpdateContent(ctx context.Context, newContent *turnip.Content) (*turnip.Content, error) {
	// todo here

	// todo: make setting more dynamic instead of setting everything
	// todo set these attributes to UpdateContent
	accessDetails := &turnip.AccessDetails{}
	meta := ""

	_, err := c.db.Pool.Exec(ctx, `UPDATE public."Content"
		SET title=$1, description=$2, content=$3, tag_list=$4, access_details=$5, meta=$6
		WHERE primary_id = $7`,
		newContent.Title, newContent.Description, newContent.Content, // 1-3
		newContent.TagList, accessDetails, meta, newContent.PrimaryId, // 4-7
	)

	// todo: put NoRowErr check here!

	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
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

func NewContentsPostgres(ctx context.Context, d *PostgresDb) Contents {
	// primary: primary id
	// sort: created at
	t := contentsPostgresImpl{
		db:        d,
		tableName: "Content",
	}

	SetupTable(ctx, d, &t)

	return &t
}
