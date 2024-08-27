package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Serve() {
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	<-termSignal
	fmt.Println("\nShutting down...")
}
