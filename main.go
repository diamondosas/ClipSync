package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"clipsync/gui"
	"clipsync/internal/clipboard"
	// firewall"clipsync/internal/firewall"
	start "clipsync/internal/init"
	// "clipsync/internal/utils"
)

var Version = "dev"

func main() {
	// if err := utils.EnsureAppInPath(); err != nil {
	// 	log.Printf("Failed to ensure app is in PATH: %v", err)
	// }

	clipboard.Init()

	// Intercept CLI execution. If it returns true, we shouldn't start GUI.

	// Setup context for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Run background sync tasks in a goroutine
	go func() {
		err := start.InitServices(ctx)
		if err != nil && err != context.Canceled {
			log.Printf("Background sync stopped: %v", err)
		}
	}()
	
	// Start the GUI (blocking call)
	gui.StartGUI()
}
