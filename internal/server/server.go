package server

import (
	"context"
	"github.com/TurnipXenon/turnip/internal/storage"
	"log"

	"github.com/TurnipXenon/turnip/internal/models"
)

type Server struct {
	Storage  storage.Storage // todo: fix
	Users    storage.Users
	Tokens   storage.Tokens
	Contents storage.Contents
	db       *storage.PostgresDb
}

// InitializeServer remember to defer cleanup!
func InitializeServer(ctx context.Context, flags models.RunFlags) *Server {
	s := Server{}

	// region db
	s.db = storage.NewPostgresDatabase(ctx, flags)

	if flags.IsLocal {
		s.Storage = storage.NewStorageLocal()
	} else {
		// todo(turnip): implement for deployment
		log.Fatalf("TODO: Unimplemented")
	}

	// todo
	s.Users = storage.NewUsersPostgres(ctx, s.db)
	// todo(turnip)
	s.Tokens = storage.NewTokensPostgres(ctx, s.db)
	// todo(turnip)
	s.Contents = storage.NewContentsPostgres(ctx, s.db)
	// endregion db

	return &s
}

func (s *Server) Cleanup(ctx context.Context) {
	if s.db != nil {
		s.db.DeferredClose(ctx)
	}
}
