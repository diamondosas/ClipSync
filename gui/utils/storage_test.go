// ai-generated
package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoragePinned(t *testing.T) {
	testItems := []string{"hello world", "git status", "https://github.com"}

	err := SavePinned(testItems)
	if err != nil {
		t.Fatalf("SavePinned failed: %v", err)
	}

	loaded := LoadPinned()
	if len(loaded) != len(testItems) {
		t.Fatalf("Expected %d items, got %d", len(testItems), len(loaded))
	}

	for i, item := range testItems {
		if loaded[i] != item {
			t.Errorf("Expected item %d to be %q, got %q", i, item, loaded[i])
		}
	}

	// Verify file location
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir() failed: %v", err)
	}
	expectedFile := filepath.Join(configDir, "clipsync", "pinned.json")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected pinned.json to exist at %s", expectedFile)
	}
}
