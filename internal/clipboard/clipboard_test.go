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

	var Outputch = make(chan []byte)

	go func(){
		data := clipboard.WatchClipboard(ctx)	
		Outputch <-data
	}()
	
	clipboard.WriteClipboard(want)
	
	Output := <-Outputch
	if !bytes.Equal(Output, []byte(want)) {
		t.Errorf("Input: %v Output: %v", want, Output)
	}
	t.Logf("Input: %v Output: %v", want, Output)
}