package clipboard_test

import (
	"clipsync/internal/clipboard"
	"testing"
	"bytes"
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	datach := clipboard.WatchClipboard(ctx)

	clipboard.WriteClipboard(want)

	text := <-datach

	if !bytes.Equal(text, []byte(want)) {
		t.Errorf("Input: %v Output: %v", want, string(text))
	}
}