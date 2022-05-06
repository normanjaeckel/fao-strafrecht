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

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/public"
	"golang.org/x/sys/unix"
)

type Logger interface {
	Printf(format string, v ...any)
}

type LoggerWithFatal interface {
	Logger
	Fatalf(format string, v ...any)
}

type Environment interface {
	Host() string
	Port() string
}

const APIPrefix = "api"

// Run is the entry point for this module. It does some preparation and then
// starts the server.
func Run(logger LoggerWithFatal, env Environment, m *model.Model) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		onSignals(logger, cancel)
		logger.Fatalf("Aborted.")
	}()

	addr := fmt.Sprintf("%s:%s", env.Host(), env.Port())
	if err := Start(ctx, logger, m, addr); err != nil {
		return err
	}

	return nil
}

func Handler(logger Logger, m *model.Model) http.Handler {
	mux := http.NewServeMux()

	// // Model case
	p := "/" + APIPrefix + "/" + "case"
	h := NewCaseHandler(logger, m)
	mux.Handle(p+"/", http.StripPrefix(p, h))

	// Root
	mux.Handle("/", public.Files())

	return mux
}

// Start starts the server. It blocks and returns an error if the server was not shut down
// gracefully.
func Start(ctx context.Context, logger Logger, m *model.Model, addr string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: Handler(logger, m),
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
func onSignals(logger Logger, cancel context.CancelFunc) {
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
