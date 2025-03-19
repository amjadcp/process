package main

import (
	"context"
	"fmt"
	"github.com/amjadcp/process/process"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx) // Store cancel function

	// Create a buffered channel to avoid blocking
	eventsChannel := make(chan string, 100)

	// Start tracking processes
	go process.TrackProcesses(eventsChannel, 2*time.Second)

	// Forward each process event to the frontend via Wails events
	go func() {
		for {
			select {
			case event, ok := <-eventsChannel:
				if !ok {
					return // Exit if channel is closed
				}
				fmt.Println(event)
				runtime.EventsEmit(a.ctx, "process_log", event)
			case <-a.ctx.Done():
				fmt.Println("Stopping process tracking...")
				return
			}
		}
	}()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Perform any post-startup setup here
}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// Optionally ask the user for confirmation before closing
	// return true to prevent closing
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	fmt.Println("Shutting down application...")

	// Stop the process tracking goroutine
	if a.cancel != nil {
		a.cancel()
	}
}
