package root

import (
	"clipsync/internal"
	"clipsync/internal/clipboard"
	"clipsync/internal/network"
	"context"
	"testing"
	"time"
)

func TestStartClipSync_CleanShutdown(t *testing.T) {
	// Temporarily switch to a test port to avoid conflicting with standard 9999
	// Change testPort if another port is needed
	testPort := "19997"
	oldPort := internal.PORT
	internal.PORT = testPort
	defer func() { internal.PORT = oldPort }()

	// Change cancel timeout to test longer or shorter runtime before shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err := StartClipSync(ctx)
	// StartClipSync should exit cleanly when context is canceled
	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		t.Errorf("StartClipSync returned unexpected error on shutdown: %v", err)
	}
}

func TestStartClipSync_ReceivePeerClip(t *testing.T) {
	testPort := "19996"
	oldPort := internal.PORT
	internal.PORT = testPort
	defer func() { internal.PORT = oldPort }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run StartClipSync in background
	syncErr := make(chan error, 1)
	go func() {
		syncErr <- StartClipSync(ctx)
	}()

	// Wait for goroutines to initialize
	time.Sleep(50 * time.Millisecond)

	// Simulate incoming clip from peer network channel
	// Change syncText to test other incoming payloads
	syncText := "Hello from peer device"
	network.IncomingClips <- []byte(syncText)

	// Allow receiver goroutine to process clip
	time.Sleep(80 * time.Millisecond)

	// Verify clipboard was updated with received data
	copied := clipboard.CopyClipboard(ctx)
	if copied != syncText {
		t.Errorf("CopyClipboard() = %q; want %q", copied, syncText)
	}

	// Verify internal history updated
	internal.ClipHistoryMu.Lock()
	hasText := len(internal.ClipHistory) > 0 && internal.ClipHistory[0] == syncText
	internal.ClipHistoryMu.Unlock()
	if !hasText {
		t.Errorf("expected ClipHistory to contain %q", syncText)
	}

	cancel()
	<-syncErr
}
