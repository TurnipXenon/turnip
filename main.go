package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"

	"github.com/TurnipXenon/turnip/internal/api"
	"github.com/TurnipXenon/turnip/internal/models"
	"github.com/TurnipXenon/turnip/internal/server"
)

func main() {
	// parse flags
	flags := models.RunFlags{}
	flag.IntVar(&flags.Port, "port", 8000, "port number to serve the turnip")
	flag.BoolVar(&flags.IsLocal, "is-local", true, "determines whether to use local services or not")
	flag.Parse()

	// load environment
	// todo: conditionally load
	err := godotenv.Load("configs/local.env")
	if err != nil {
		log.Fatalf("Failed loading local environment. Err: %s", err)
	}

	// todo: set up connections to other services like db
	// ctx := context.Background()
	s := server.InitializeServer(flags)

	// run serve mux or router
	api.RunServeMux(s, flags)
}
