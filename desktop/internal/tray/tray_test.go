package tray

import (
	"bytes"
	"testing"
)

func TestTrayIcons_Embedded(t *testing.T) {
	// Verify PNG icon is embedded and has valid PNG magic bytes
	if len(trayIcon_png) == 0 {
		t.Fatal("expected embedded trayIcon_png to not be empty")
	}
	pngHeader := []byte("\x89PNG\r\n\x1a\n")
	if !bytes.HasPrefix(trayIcon_png, pngHeader) {
		t.Errorf("trayIcon_png does not have standard PNG header")
	}

	// Verify ICO icon is embedded and has valid ICO magic header
	if len(trayIcon_ico) == 0 {
		t.Fatal("expected embedded trayIcon_ico to not be empty")
	}
	icoHeader := []byte{0x00, 0x00, 0x01, 0x00}
	if !bytes.HasPrefix(trayIcon_ico, icoHeader) {
		t.Errorf("trayIcon_ico does not have standard ICO header")
	}
}

// func TestOnTrayExit_Subprocess(t *testing.T) {
// 	// Subprocess execution pattern to verify OnTrayExit calls cancel and exits with code 0
// 	if os.Getenv("BE_TRAY_EXIT_SUBPROCESS") == "1" {
// 		calledCancel := false
// 		cancel := func() {
// 			calledCancel = true
// 		}
// 		OnTrayExit(cancel)
// 		return
// 	}

// 	cmd := exec.Command(os.Args[0], "-test.run=TestOnTrayExit_Subprocess")
// 	cmd.Env = append(os.Environ(), "BE_TRAY_EXIT_SUBPROCESS=1")
// 	err := cmd.Run()
// 	if err != nil {
// 		t.Fatalf("expected subprocess to exit cleanly with 0, got error: %v", err)
// 	}
// }
