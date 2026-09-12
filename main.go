package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"clipsync/gui"
	"clipsync/internal/root"
	"clipsync/internal/tray"

	"github.com/getlantern/systray"
)

func main() {
	// Setup context for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Run background Services in a goroutine
	go func() {
		log.Println("Starting background sync services...")
		err := root.StartClipSync(ctx)
		if err != nil && err != context.Canceled {
			log.Printf(" Background sync stopped with error: %v", err)
		} else {
			log.Println("Background sync stopped cleanly")
		}
	}()

	systray.Register(func() { tray.OnTrayReady(cancel) }, tray.OnTrayExit)
	
	log.Println("[Main] Launching GUI...")
	gui.StartGUI()
}
