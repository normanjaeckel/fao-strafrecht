package main

import (
	"log"
	"os"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/db"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/env"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
)

func main() {
	// Provide dependencies. See pkg/deps for more information.

	// Logger
	logger := log.Default()

	// Environment
	environment, err := env.Parse(os.Getenv)
	if err != nil {
		logger.Fatalf("Error: parsing environment: %v", err)
	}

	// Database
	f, err := os.OpenFile(
		environment.DBFilename(),
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0600,
	)
	if err != nil {
		logger.Fatalf("Error: opening database file: %v", err)
	}
	database, err := db.New(f)
	if err != nil {
		logger.Fatalf("Error: loading database file: %v", err)
	}

	// Start everything.
	if err := srv.Run(logger, environment, database); err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
