package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/TurnipXenon/turnip/internal/models"
	"github.com/TurnipXenon/turnip/internal/storage/migration"
	"github.com/TurnipXenon/turnip/internal/util"
)

const (
	ifUserTableExists = `SELECT EXISTS(
               SELECT
               FROM information_schema.tables
               WHERE table_schema = 'public'
                 AND table_name = '%s'
           );`
)

type PostgresDb struct {
	Pool *pgxpool.Pool
}

type GenericPostgresTable interface {
	GetTableName() string
	GetMigrationSequence() []migration.Migration
}

// NewPostgresDatabase remember to defer DeferredClose!!!
// todo: improve documentation
func NewPostgresDatabase(ctx context.Context, flags models.RunFlags) *PostgresDb {
	p := PostgresDb{}
	var err error
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	config, err := pgxpool.ParseConfig(flags.PostgresConnection)
	if err != nil {
		log.Fatalf("Unable to create connection config: %v\n", err)
	}
	p.Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v\n", err)
	}

	return &p
}

// todo safely close
func (p *PostgresDb) DeferredClose(ctx context.Context) {
	p.Pool.Close()
}

func IfTableExistsQuery(tableName string) string {
	return fmt.Sprintf(ifUserTableExists, tableName)
}

func SetupTable(ctx context.Context, d *PostgresDb, t GenericPostgresTable) {
	// todo: refactor to be a better, the side-effect is not good
	row := d.Pool.QueryRow(ctx, IfTableExistsQuery(t.GetTableName()))
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		util.LogDetailedError(err)
		log.Fatalf("failed to check if table exists: %v", err)
	}

	if !exists {
		for _, m := range t.GetMigrationSequence() {
			m.Migrate(ctx, d.Pool)
		}
	}
}
