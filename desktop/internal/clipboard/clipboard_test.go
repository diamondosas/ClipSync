package clipboard

import (
	"clipsync/internal"
	"context"
	"testing"
	"time"
)

func TestClipboard_ReadWrite(t *testing.T) {
	ctx := context.Background()
	// Change sampleText to test different clipboard strings
	sampleText := "Player1 high score: 9999"

	WriteClipboard(ctx, sampleText)

	got := CopyClipboard(ctx)
	if got != sampleText {
		t.Errorf("CopyClipboard() = %q; want %q", got, sampleText)
	}
}

func TestClipboard_UnicodeAndSpecialChars(t *testing.T) {
	ctx := context.Background()
	// Change unicodeText to test other multi-byte or emoji payloads
	unicodeText := "ClipSync 🚀 剪贴板 \t\n 123"

	WriteClipboard(ctx, unicodeText)

	got := CopyClipboard(ctx)
	if got != unicodeText {
		t.Errorf("CopyClipboard() = %q; want %q", got, unicodeText)
	}
}

func TestClipboard_EmptyString(t *testing.T) {
	ctx := context.Background()
	// Change emptyPayload if testing whitespace or blank strings
	emptyPayload := ""

	WriteClipboard(ctx, emptyPayload)

	got := CopyClipboard(ctx)
	if got != emptyPayload {
		t.Errorf("CopyClipboard() = %q; want %q", got, emptyPayload)
	}
}

func TestWatchClipboard_IgnoreNetworkClip(t *testing.T) {
	// Change timeout duration here if testing longer watcher lifetimes
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Simulate incoming network clip so WatchClipboard ignores it
	ignoredPayload := []byte("sync-from-laptop")
	internal.LastRecvClipMu.Lock()
	internal.LastRecvClip = ignoredPayload
	internal.LastRecvClipMu.Unlock()

	WriteClipboard(ctx, string(ignoredPayload))

	// Write a new real text shortly after
	realPayload := "New game invite code"
	go func() {
		time.Sleep(100 * time.Millisecond)
		WriteClipboard(ctx, realPayload)
	}()

	got := WatchClipboard(ctx)
	res := string(<-got)
	if res != realPayload {
		t.Errorf("WatchClipboard() = %q; want %q", res, realPayload)
	}
}

func TestWatchClipboard_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	out := WatchClipboard(ctx)

	// Cancel context immediately
	cancel()

	// Ensure no hang or panic after context closure
	select {
	case <-out:
	case <-time.After(50 * time.Millisecond):
	}
}
