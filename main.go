package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amjadcp/process/process"
)

func main() {
	events := make(chan string) // Channel to receive process events

	// Start tracking processes in a separate goroutine
	go process.TrackProcesses(events)

	// Handle Ctrl+C (SIGINT) for graceful exit
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("🚀 Process monitoring started. Waiting for changes...\n")

	for {
		select {
		case event := <-events:
			fmt.Println(event)
		case <-stop:
			fmt.Println("\nExiting gracefully...")
			return
		}
	}
}
