package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
	"golang.org/x/sys/unix"
)

const defaultAddr = ":8000"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, unix.SIGINT, unix.SIGTERM)
		<-sig
		cancel()
		<-sig
		log.Fatal("Abborted")
	}()
	if err := srv.Start(ctx, defaultAddr); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
