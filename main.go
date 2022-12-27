package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/TurnipXenon/turnip/internal/api"
	"github.com/TurnipXenon/turnip/internal/models"
	"github.com/TurnipXenon/turnip/internal/server"
)

func main() {
	flags := models.RunFlags{}

	// override with environment
	flags.PostgresConnection = os.Getenv("PGCONN")
	if flags.PostgresConnection == "" {
		// local setup
		flags.PostgresConnection = "postgresql://turnipservice:password@localhost:5432/turnip"
	}

	default_port_str := os.Getenv("PORT")
	var port int
	if default_port_str == "" {
		port = 80
	} else {
		port, _ = strconv.Atoi(default_port_str)
		// todo: handle error
	}

	// parse flags
	flag.IntVar(&flags.Port, "port", port, "port number to serve the turnip")
	flag.BoolVar(&flags.IsLocal, "is-local", true, "determines whether to use local services or not")
	flag.Parse()

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
	api.RunServeMux(s, flags)
}
