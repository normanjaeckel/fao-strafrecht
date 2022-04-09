package srv_test

import (
	"context"
	"log"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
)

func TestStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error, 1)

	logger := log.Default()

	go func() {
		ch <- srv.Start(ctx, logger, ":8000")
	}()
	cancel()

	srvErr := <-ch
	if srvErr != nil {
		t.Fatalf("got error from closed server: %v", srvErr)
	}

}
