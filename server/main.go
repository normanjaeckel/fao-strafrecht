package main

import (
	"log"
	"os"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
)

func main() {
	// Provide dependencies. See pkg/deps for more information.
	logger := log.Default()
	getEnvFunc := os.Getenv

	// Start everything.
	if err := srv.Run(logger, getEnvFunc); err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
