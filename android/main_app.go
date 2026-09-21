package main

import (
	"clipsync-android/gui"
	"clipsync-android/internal/events"
	"clipsync-android/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main_app() {
	// Setup context for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize decoupled event bus and service engine
	bus := events.NewEventBus()
	svc := service.NewClipSyncService(bus, service.Config{})

	// Run background networking service
	go func() {
		log.Println("[Main] Starting ClipSync background network service...")
		if err := svc.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("[Main] ClipSync service error: %v", err)
		} else {
			log.Println("[Main] ClipSync service stopped cleanly")
		}
	}()

	// Launch Gio Mobile GUI
	log.Println("[Main] Starting ClipSync Gio UI...")
	gui.StartGUI(svc)
}
