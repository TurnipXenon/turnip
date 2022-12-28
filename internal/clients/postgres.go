package clients

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/TurnipXenon/turnip/internal/models"
)

type PostgresDb struct {
	Pool *pgxpool.Pool
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
