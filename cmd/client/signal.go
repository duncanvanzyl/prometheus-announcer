package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// ExitOnSignal calls the provided CancelFunc on receipt of SIGINT or SIGTERM
func ExitOnSignal(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("exit on signal: %s=%v", "signal", sig)
	cancel()
}
