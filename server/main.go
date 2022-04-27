package main

import (
	"log"
	"os"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/env"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/eventstore"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
)

func main() {
	// Logger
	logger := log.Default()

	// Environment
	environment, err := env.Parse(os.Getenv)
	if err != nil {
		logger.Fatalf("Error: parsing environment: %v", err)
	}

	// Datastore
	es, close, err := eventstore.New(logger, environment.DSFilename())
	if err != nil {
		logger.Fatalf("Error: loading eventstore: %v", err)
	}
	defer close()

	// Start everything.
	if err := srv.Run(logger, environment, es); err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
