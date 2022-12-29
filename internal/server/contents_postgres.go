package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"

	"github.com/TurnipXenon/turnip/internal/clients"
	"github.com/TurnipXenon/turnip/internal/server/sql/migration"
	"github.com/TurnipXenon/turnip/internal/util"
)

type contentsPostgresImpl struct {
	db        *clients.PostgresDb
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

func (c *contentsPostgresImpl) CreateContent(ctx context.Context, request *turnip.CreateContentRequest, user *turnip.User) (*turnip.Content, error) {
	// todo: require some fields!

	// create uuid
	// very unlikely to collide, right?
	content := turnip.Content{
		Title:         request.Title,
		Description:   request.Description,
		Content:       request.Content,
		Media:         request.Media,
		TagList:       request.TagList,
		AccessDetails: request.AccessDetails,
		Meta:          request.Meta,
		PrimaryId:     uuid.New().String(),
		CreatedAt:     timestamppb.Now(),
		AuthorId:      user.PrimaryId,
	}
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

	return &content, nil
}

func (c *contentsPostgresImpl) GetContentById(ctx context.Context, primary string) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c *contentsPostgresImpl) GetAllContent(ctx context.Context) ([]*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c *contentsPostgresImpl) GetContentByTag(ctx context.Context, tag string) ([]*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c *contentsPostgresImpl) UpdateContent(ctx context.Context, new *turnip.Content) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c *contentsPostgresImpl) DeleteContentById(ctx context.Context, primary string) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func NewContentsPostgres(ctx context.Context, d *clients.PostgresDb) Contents {
	// primary: primary id
	// sort: created at
	t := contentsPostgresImpl{
		db:        d,
		tableName: "Content",
	}

	clients.SetupTable(ctx, d, &t)

	return &t
}
