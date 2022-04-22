package main

import (
	"log"
	"os"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/db"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
)

func main() {
	// Provide dependencies. See pkg/deps for more information.
	logger := log.Default()
	db := db.New()
	getEnvFunc := os.Getenv

	// Start everything.
	if err := srv.Run(logger, db, getEnvFunc); err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
