package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amjadcp/process/process"
)

const pollInterval = 2 * time.Second

func main() {
	events := make(chan string)

	// Start tracking processes in a separate goroutine.
	go process.TrackProcesses(events, pollInterval)

	// Handle graceful shutdown (Ctrl+C or SIGTERM).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("🚀 Process monitoring started. Waiting for changes...")

	for {
		select {
		case event := <-events:
			fmt.Println(event)
		case <-stop:
			fmt.Println("\nExiting gracefully...")
			close(events)
			// Give some time for any pending asynchronous operations.
			time.Sleep(500 * time.Millisecond)
			return
		}
	}
}
