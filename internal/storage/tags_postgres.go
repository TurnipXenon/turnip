package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"

	"github.com/TurnipXenon/turnip/internal/storage/migration"
	"github.com/TurnipXenon/turnip/internal/util"
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
	// get old tags
	oldTagList, err := t.GetTagsByContent(ctx, content)
	if err != nil {
		util.LogDetailedError(err)
		return err
	}

	// turn old and new to maps or psuedo-sets
	newTagMap := map[string]bool{}
	for _, s := range content.TagList {
		newTagMap[strings.ToLower(s)] = true
	}
	oldTagMap := map[string]bool{}
	for _, s := range oldTagList {
		oldTagMap[s] = true // we trust this one to be already normalized
	}

	// if old tag not in new tag, delete
	var tagsToDeleteList []string
	for _, s := range oldTagList {
		if !newTagMap[s] {
			tagsToDeleteList = append(tagsToDeleteList, fmt.Sprintf("'%s'", s))
		}
	}
	if len(tagsToDeleteList) > 0 {
		tagsToDeleteStr := strings.Join(tagsToDeleteList, ", ")
		_, err = t.db.Pool.Exec(ctx,
			fmt.Sprintf(`DELETE FROM "%s" WHERE tag IN (%s) AND content_id=$1`,
				t.tableName, tagsToDeleteStr),
			content.PrimaryId)
	}

	// if new tag not in old tag, create
	for _, s := range content.TagList {
		normalized := strings.ToLower(s)
		if !oldTagMap[normalized] {
			// todo: create
			_, err = t.db.Pool.Exec(ctx,
				`INSERT INTO "Tag" (tag, content_id, created_at) 
					VALUES ($1, $2, $3)`,
				normalized, content.PrimaryId, content.CreatedAt.AsTime().Format(time.RFC3339))
			if err != nil {
				util.LogDetailedError(err)
			}
		}
	}

	return nil
}

func (t *tagsPostgresImpl) GetTagsByContent(ctx context.Context, content *turnip.Content) ([]string, error) {
	rows, _ := t.db.Pool.Query(ctx,
		`SELECT tag FROM "Tag" WHERE content_id=$1`, content.PrimaryId)
	oldTagList, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var n string
		err := row.Scan(&n)
		return n, err
	})
	if err != nil {
		util.LogDetailedError(err)
		return []string{}, util.WrapErrorWithDetails(err)
	}

	return oldTagList, nil
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
