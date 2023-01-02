package storage

import (
	"context"
	"errors"

	"github.com/TurnipXenon/turnip/internal/storage/migration"
	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

type tagsPostgresImpl struct {
	db        *PostgresDb
	tableName string
}

func (t *tagsPostgresImpl) DeleteTags(ctx context.Context, primaryId string) error {
	//TODO implement me
	panic("implement me")
}

func (t *tagsPostgresImpl) UpdateTags(ctx context.Context, content *turnip.Content) error {
	//TODO implement me
	return errors.New("unimplemented")
}

func (t *tagsPostgresImpl) GetTagsByContent(ctx context.Context, content *turnip.Content) ([]string, error) {
	//TODO implement me
	return []string{}, errors.New("unimplemented")
}

func (t *tagsPostgresImpl) GetContentIdsByTag(ctx context.Context, tagList []string) ([]string, error) {
	//TODO implement me
	return []string{}, errors.New("unimplemented")
}

func (t *tagsPostgresImpl) GetTableName() string {
	return t.tableName
}

func (t *tagsPostgresImpl) GetMigrationSequence() []migration.Migration {
	return []migration.Migration{
		migration.NewGenericMigration(migration.MigrateTags0001),
	}
}

func NewTagsPostgres(ctx context.Context, d *PostgresDb) Tags {
	t := tagsPostgresImpl{
		db:        d,
		tableName: "Tag",
	}

	SetupTable(ctx, d, &t)

	return &t
}
