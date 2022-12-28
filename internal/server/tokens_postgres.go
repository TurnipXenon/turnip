// storage is an abstraction to s3 buckets

package server

import (
	"context"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"

	"github.com/TurnipXenon/turnip/internal/clients"
	"github.com/TurnipXenon/turnip/internal/server/sql/migration"
)

type tokensPostgresImpl struct {
	db          *clients.PostgresDb
	dbTableName string
	// todo: global secondary index
}

func (t *tokensPostgresImpl) GetTableName() string {
	return t.dbTableName
}

func (t *tokensPostgresImpl) GetMigrationSequence() []migration.Migration {
	return []migration.Migration{
		migration.NewGenericMigration(migration.MigrateToken0001),
	}
}

func (t tokensPostgresImpl) GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*turnip.Token, error) {
	//TODO implement me
	panic("implement me")
}

func (t tokensPostgresImpl) GetToken(token string) (*turnip.Token, error) {
	//TODO implement me
	panic("implement me")
}

func NewTokensPostgres(ctx context.Context, d *clients.PostgresDb) Tokens {
	t := tokensPostgresImpl{
		db:          d,
		dbTableName: "Tokens",
	}

	clients.SetupTable(ctx, d, &t)

	return &t
}
