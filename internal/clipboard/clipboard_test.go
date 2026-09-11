//@ai-generated
package clipboard

import (
	"context"
	"testing"
	"time"
	"clipsync/internal/network"
)

func TestClipboard_ReadWrite(t *testing.T) {
	ctx := context.Background()
	sampleText := "Player1 high score: 9999"

	WriteClipboard(ctx, sampleText)

	got := CopyClipboard(ctx)
	if got != sampleText {
		t.Errorf("CopyClipboard() = %q; want %q", got, sampleText)
	}
}

func TestWatchClipboard_ContextCancel(t *testing.T) {
	// Cancel quickly to ensure WatchClipboard unblocks
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	data := WatchClipboard(ctx)
	if data != nil {
		t.Errorf("WatchClipboard() with canceled context = %v; want nil", data)
	}
}

func TestWatchClipboard_IgnoreNetworkClip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Simulate incoming network clip so WatchClipboard ignores it
	ignoredPayload := []byte("sync-from-laptop")
	network.LastRecievedClip = ignoredPayload

	// Write the ignored text
	WriteClipboard(ctx, string(ignoredPayload))

	// Write a new real text shortly after
	realPayload := "New game invite code"
	go func() {
		time.Sleep(100 * time.Millisecond)
		WriteClipboard(ctx, realPayload)
	}()

	got := WatchClipboard(ctx)
	if string(<-got) != realPayload {
		t.Errorf("WatchClipboard() = %q; want %q", string(<-got), realPayload)
	}
}