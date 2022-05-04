package main

import (
	"log"
	"os"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/env"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/eventstore"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
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

	// Eventstore
	es, close, err := eventstore.New(logger, environment.DSFilename())
	if err != nil {
		logger.Fatalf("Error: loading eventstore: %v", err)
	}
	defer close()

	model, err := model.New(es)
	if err != nil {
		logger.Fatalf("Error: loading model: %v", err)
	}

	// Start everything.
	if err := srv.Run(logger, environment, model); err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
