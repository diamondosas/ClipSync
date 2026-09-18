// @ai-generated
package clipboard

import (
	"clipsync/internal"
	"context"
	"testing"
	"time"
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

func TestWatchClipboard_IgnoreNetworkClip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Simulate incoming network clip so WatchClipboard ignores it
	ignoredPayload := []byte("sync-from-laptop")
	internal.LastRecvClipMu.Lock()
	internal.LastRecvClip = ignoredPayload
	internal.LastRecvClipMu.Unlock()
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
