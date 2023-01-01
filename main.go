package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/TurnipXenon/turnip/internal/models"
	"github.com/TurnipXenon/turnip/internal/server"
)

func main() {
	flags := models.InitializeFlags()

	// load environment
	// todo: conditionally load
	err := godotenv.Load("configs/local.env")
	if err != nil {
		log.Fatalf("Failed loading local environment. Err: %s", err)
	}

	// todo: set up connections to other services like db
	ctx := context.Background()
	s := server.InitializeServer(ctx, flags)
	defer s.Cleanup(ctx)

	// run serve mux or router
	server.RunServeMux(s, flags)
}
