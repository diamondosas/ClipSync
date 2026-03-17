package clipboard_test

import (
	"clipsync/internal/clipboard"
	"testing"
	"bytes"
	"context"
	"os/signal"
	"os"
	"syscall"
)

func TestReadWrite(t *testing.T) {
	want := "Testing is taking place..."
	clipboard.WriteClipboard(want)
	output := clipboard.CopyClipboard()

	if want != output {
		t.Errorf("Input: %v Output : %v", want, output)
	}

}

func TestWatch(t *testing.T) {
	want := "Tester"
	
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	var text = make(chan []byte)
	go clipboard.WatchClipboard(ctx, text)
	clipboard.WriteClipboard(want)
	output := <-text

	if !bytes.Equal(output, []byte(want)) {
		t.Errorf("Input: %v Output: %v", want, output)
	}
	t.Logf("Input: %v Output: %v", want, output)
}