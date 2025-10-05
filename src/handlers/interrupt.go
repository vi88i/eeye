package handlers

import (
	"os"
	"os/signal"
	"syscall"
)

// GetInterruptHandlerChannel returns a channel that can be used to capture keyboard interrupts
func GetInterruptHandlerChannel() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	return quit
}
