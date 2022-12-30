package storage

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/TurnipXenon/turnip/internal/storage/migration"
	"github.com/TurnipXenon/turnip/internal/util"
)

const (
	metadataCanUsersBeMadeKey = "can_users_be_made"
)

type metadataPostgresImpl struct {
	db        *PostgresDb
	tableName string
}

func (m *metadataPostgresImpl) CanUsersBeMade(ctx context.Context) (bool, error) {
	// todo here
	row := m.db.Pool.QueryRow(
		ctx,
		`SELECT value FROM "Metadata" WHERE key=$1`,
		metadataCanUsersBeMadeKey)
	canUsersBeMade := false
	var canUsersBeMadeStr string
	err := row.Scan(&canUsersBeMadeStr)
	if err == pgx.ErrNoRows {
		// todo: document behavior: if row does not exist, make one! if succ, return true, if not, return false
		canUsersBeMade, err = m.SetCanUsersBeMade(ctx, true)
	} else if err != nil {
		// todo: document behavior
		util.LogDetailedError(err)
		return false, util.WrapErrorWithDetails(err)
	} else {
		canUsersBeMade, err = strconv.ParseBool(canUsersBeMadeStr)
		if err != nil {
			util.LogDetailedError(err)
			return false, util.WrapErrorWithDetails(err)
		}
	}
	return canUsersBeMade, nil
}

func (m *metadataPostgresImpl) SetCanUsersBeMade(ctx context.Context, canBeMade bool) (bool, error) {
	_, err := m.db.Pool.Exec(ctx,
		`INSERT INTO "Metadata"
    (key, value) 
	VALUES ($1, $2)`,
		metadataCanUsersBeMadeKey, strconv.FormatBool(canBeMade),
	)

	if err != nil {
		util.LogDetailedError(err)
		return false, util.WrapErrorWithDetails(err)
	}

	return canBeMade, nil
}

func (m *metadataPostgresImpl) GetTableName() string {
	return m.tableName
}

func (m *metadataPostgresImpl) GetMigrationSequence() []migration.Migration {
	return []migration.Migration{
		migration.NewGenericMigration(migration.MigrateMetadata0001),
	}
}

func NewMetadataPostgres(ctx context.Context, d *PostgresDb) Metadata {
	t := metadataPostgresImpl{
		db:        d,
		tableName: "Metadata",
	}

	SetupTable(ctx, d, &t)

	return &t
}
