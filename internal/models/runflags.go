package models

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/TurnipXenon/turnip/internal/util"
)

type RunFlags struct {
	Port               int
	IsLocal            bool
	PostgresConnection string
	CorsAllowList      []string
}

func InitializeFlags() *RunFlags {
	flags := RunFlags{}

	// override with environment
	flags.PostgresConnection = os.Getenv("DATABASE_URL")
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

	// environment variable exclusives
	corsAllowListStr := os.Getenv("CORS_ALLOWLIST")
	err := json.Unmarshal([]byte(corsAllowListStr), &flags.CorsAllowList)
	if err != nil {
		util.LogDetailedError(err)
	}
	if flags.CorsAllowList == nil {
		flags.CorsAllowList = []string{}
	}

	return &flags
}
