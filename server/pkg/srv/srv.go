/*
Package srv provides the HTTP server for the backend.
*/
package srv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/deps"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/env"
	"golang.org/x/sys/unix"
)

// Run is the entry point for this module. It does some preparation and then
// starts the server.
func Run(logger deps.Logger, getEnvFunc deps.GetEnvFunc) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		onSignals(logger, cancel)
		logger.Fatalf("Aborted.")
	}()

	e, err := env.Parse(getEnvFunc)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)
	if err := Start(ctx, logger, addr); err != nil {
		return err
	}

	return nil
}

// Start starts the server. It blocks and returns an error if the server was not shut down
// gracefully.
func Start(ctx context.Context, logger deps.Logger, addr string) error {
	s := &http.Server{
		Addr: addr,
	}

	go func() {
		<-ctx.Done()
		logger.Printf("Server is shuting down")
		if err := s.Shutdown(context.Background()); err != nil {
			logger.Printf("Error: Shutting down server: %v", err)
		}
	}()

	logger.Printf("Server starts and listens on %q", addr)
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		return fmt.Errorf("server exited: %w", err)
	}
	logger.Printf("Server is down")

	return nil
}

// onSignals blocks until the operating system sends SIGTERM or SIGINT. On
// incomming signal, it calls the cancel function and blocks again until SIGINT
// comes in. Use it in a goroutine and call os.Exit() with non zero exit code to
// abort the process.
func onSignals(logger deps.Logger, cancel context.CancelFunc) {
	msg := "Received operating system signal: %s"

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, unix.SIGTERM)
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, unix.SIGINT)

	select {
	case s := <-sigInt:
		logger.Printf(msg, s.String())
		cancel()
	case s := <-sigTerm:
		logger.Printf(msg, s.String())
		cancel()
	}

	s := <-sigInt
	logger.Printf(msg, s.String())
}
