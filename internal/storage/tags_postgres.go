package storage

import (
	"context"
	"fmt"
	"sort"
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

func (t *tagsPostgresImpl) DeleteTagsByContentId(ctx context.Context, primaryId string) error {
	// todo: consider deleting in the future
	_, err := t.db.Pool.Query(ctx,
		`SELECT FROM "Tag" WHERE content_id=$1`, primaryId)
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}
	return nil
}

func stringListToSqlInArgument(valueList []string) string {
	if len(valueList) == 0 {
		return ""
	}

	var l []string
	for _, s := range valueList {
		l = append(l, fmt.Sprintf("'%s'", strings.ToLower(s)))
	}
	return strings.Join(l, ", ")
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

func (t *tagsPostgresImpl) GetContentIdsByTagInclusive(ctx context.Context, tagList []string) ([]string, error) {
	if len(tagList) == 0 {
		return nil, nil
	}

	inParam := stringListToSqlInArgument(tagList)
	rows, _ := t.db.Pool.Query(ctx,
		fmt.Sprintf(`SELECT content_id FROM "Tag" WHERE tag in (%s)`, inParam))
	contentIdList, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var n string
		err := row.Scan(&n)
		return n, err
	})
	if err != nil {
		util.LogDetailedError(err)
		return []string{}, util.WrapErrorWithDetails(err)
	}

	return contentIdList, nil
}

func (t *tagsPostgresImpl) GetContentIdsByTagStrict(ctx context.Context, tagList []string) ([]string, error) {
	// todo here
	if len(tagList) == 0 {
		return nil, nil
	}

	// todo: order strings
	var l []string
	for _, s := range tagList {
		l = append(l, fmt.Sprintf("'%s'", strings.ToLower(s)))
	}
	sort.Strings(l)
	inParam := strings.Join(l, ", ")
	query := fmt.Sprintf(`SELECT content_id
FROM "Tag"
GROUP BY content_id
HAVING array_agg(tag order by tag) = array [%s]`, inParam)
	rows, _ := t.db.Pool.Query(ctx, query)
	contentIdList, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var n string
		err := row.Scan(&n)
		return n, err
	})
	if err != nil {
		util.LogDetailedError(err)
		return []string{}, util.WrapErrorWithDetails(err)
	}

	return contentIdList, nil
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
