package clients

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/TurnipXenon/turnip/internal/models"
)

type PostgresDb struct {
	Conn *pgx.Conn
}

// NewPostgresDatabase remember to defer DeferredClose!!!
// todo: improve documentation
func NewPostgresDatabase(ctx context.Context, flags models.RunFlags) *PostgresDb {
	p := PostgresDb{}
	var err error
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	p.Conn, err = pgx.Connect(ctx, flags.PostgresConnection)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v\n", err)
	}

	return &p
}

// todo safely close
func (p *PostgresDb) DeferredClose(ctx context.Context) {
	p.Conn.Close(ctx)
}
