package server

import (
	"log"

	"github.com/TurnipXenon/Turnip/internal/models"
)

type Server struct {
	Storage Storage
}

func InitializeServer(flags models.RunFlags) *Server {
	s := Server{}

	if flags.IsLocal {
		s.Storage = NewStorageLocal()
	} else {
		// todo(turnip): implement for deployment
		log.Fatalf("TODO: Unimplemented")
	}

	return &s
}
