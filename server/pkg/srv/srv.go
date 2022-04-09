/*
Package srv provides the HTTP server for the backend.
*/
package srv

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Start starts the server. It returns 0 if the server was shut down gracefully.
func Start(ctx context.Context, addr string) error {
	s := &http.Server{
		Addr: addr,
	}
	go func() {
		<-ctx.Done()
		log.Printf("Server is shutting down")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("Error: Shutting down server: %v", err)
		}
	}()
	log.Printf("Server is starting and listening at %s", addr)
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		return fmt.Errorf("server exited: %w", err)
	}
	return nil
}
