package server

import (
	"context"
	"log"

	"github.com/TurnipXenon/turnip/internal/clients"
	"github.com/TurnipXenon/turnip/internal/models"
)

type Server struct {
	Storage  Storage // todo: fix
	Users    Users
	Tokens   Tokens
	Contents Contents
	db       *clients.PostgresDb
}

// InitializeServer remember to defer cleanup!
func InitializeServer(ctx context.Context, flags models.RunFlags) *Server {
	s := Server{}

	// region db
	s.db = clients.NewPostgresDatabase(ctx)

	if flags.IsLocal {
		s.Storage = NewStorageLocal()
	} else {
		// todo(turnip): implement for deployment
		log.Fatalf("TODO: Unimplemented")
	}

	// todo
	s.Users = NewUsersPostgres(ctx, s.db)
	// todo(turnip)
	s.Tokens = NewTokensPostgres(s.db)
	// todo(turnip)
	s.Contents = NewContentsPostgres(s.db)
	// endregion db

	return &s
}

func (s *Server) Cleanup(ctx context.Context) {
	if s.db != nil {
		s.db.DeferredClose(ctx)
	}
}
