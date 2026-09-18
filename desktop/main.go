package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"clipsync/gui"
	"clipsync/internal/root"
	"clipsync/internal/tray"
	"clipsync/internal/instance"

	"github.com/getlantern/systray"
)

func main() {

	// Check for an existing instance of the application
	if instance.IsInstanceExisting() {
		os.Exit(0)
	}

	// Setup context for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	defer instance.Cleanup()

	//-- open GUI when another instance is launched
	instance.StartListener(ctx, gui.ShowWindow)
	
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
	
	//Prepare Background Tray
	go func() {
		runtime.LockOSThread()
		systray.Run(
			func() { tray.OnTrayReady(cancel, gui.ShowWindow)},
			func() { tray.OnTrayExit(cancel) },
		)
	}()
	log.Println("[Main] Launching GUI...")

	gui.StartGUI()
}
