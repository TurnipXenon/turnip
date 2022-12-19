package main

import (
	"flag"

	"github.com/TurnipXenon/Turnip/internal/api"
	"github.com/TurnipXenon/Turnip/internal/models"
)

func main() {
	// parse flags
	flags := models.RunFlags{}
	flag.IntVar(&flags.Port, "port", 8000, "port number to serve the server")
	flag.Parse()

	// ctx := context.Background()

	// todo: set up connections to other services like db

	// setup server
	api.InitializeServer(flags)
}
