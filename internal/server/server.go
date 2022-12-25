package server

import (
	"log"

	"github.com/TurnipXenon/Turnip/internal/clients"
	"github.com/TurnipXenon/Turnip/internal/models"
)

type Server struct {
	Storage  Storage
	Users    Users
	Tokens   Tokens
	Contents Contents
}

func InitializeServer(flags models.RunFlags) *Server {
	s := Server{}

	ddb := clients.NewDynamoDB()

	if flags.IsLocal {
		s.Storage = NewStorageLocal()
	} else {
		// todo(turnip): implement for deployment
		log.Fatalf("TODO: Unimplemented")
	}

	s.Users = NewUsersDynamoDB(ddb)
	s.Tokens = NewTokensDynamoDB(ddb)
	s.Contents = NewContentsDynamoDB(ddb)

	return &s
}
