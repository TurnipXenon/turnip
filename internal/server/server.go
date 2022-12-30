package server

import (
	"context"
	"fmt"
	"log"

	"github.com/TurnipXenon/turnip/internal/config"
	"github.com/TurnipXenon/turnip/internal/models"
	"github.com/TurnipXenon/turnip/internal/storage"
	"github.com/TurnipXenon/turnip/internal/util"
)

type Server struct {
	Storage  storage.Storage // todo: fix
	Users    storage.Users
	Tokens   storage.Tokens
	Contents storage.Contents
	db       *storage.PostgresDb
	Metadata storage.Metadata
}

// InitializeServer remember to defer cleanup!
func InitializeServer(ctx context.Context, flags models.RunFlags) *Server {
	s := Server{}
	s.db = storage.NewPostgresDatabase(ctx, flags)

	// region to be extracted

	// todo: extract this!
	s.Metadata = storage.NewMetadataPostgres(ctx, s.db)
	canUserBeMade, err := s.Metadata.CanUsersBeMade(ctx)
	if err != nil {
		util.LogDetailedError(err)
		fmt.Println("Will try to close CreateUser endpoint as a result")
	}

	sysConf := config.SystemConfig{CanUserBeMade: config.NewGenericSystemVariable[bool](canUserBeMade)}
	// endregion to be extracted

	// region db

	if flags.IsLocal {
		s.Storage = storage.NewStorageLocal()
	} else {
		// todo(turnip): implement for deployment
		log.Fatalf("TODO: Unimplemented")
	}

	// todo
	s.Users = storage.NewUsersPostgres(ctx, s.db, sysConf)
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
