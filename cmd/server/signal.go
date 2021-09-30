package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ExitOnSignal calls the provided CancelFunc on receipt of SIGINT or SIGTERM
func ExitOnSignal(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	logger.Info("exit on signal", "signal", sig)
	cancel()
}
