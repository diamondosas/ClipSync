package clipboard_test

import (
	"testing"
	"clipsync/internal/clipboard"
)

func TestReadWrite(t *testing.T) {
	want := "Testing is taking place..."
	clipboard.WriteClipboard(want)
	output := clipboard.CopyClipboard()

	if want != output {
		t.Errorf("Input: %v Output : %v", want, output)
	}

}

// func TestChanged(t *testing.T){
// 	want := "Tester"
// 	go clipboard.ChangedClipbord(t.Context())
	
// }