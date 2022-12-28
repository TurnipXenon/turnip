package migration

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/TurnipXenon/turnip/internal/util"
)

type Migration interface {
	Migrate(ctx context.Context, pool *pgxpool.Pool)
}

type genericMigration struct {
	sql string
}

func (g *genericMigration) Migrate(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, g.sql)
	if err != nil {
		util.LogDetailedError(err)
		log.Fatal(err)
	}
}

func NewGenericMigration(migrationScript string) Migration {
	m := genericMigration{
		sql: migrationScript,
	}
	return &m
}
