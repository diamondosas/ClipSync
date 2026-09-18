package storage

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestStorage_SaveAndLoadPinned(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "clipsync-storage-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewStorageWithDir(tempDir)

	// Initially empty
	loaded := store.LoadPinned()
	if len(loaded) != 0 {
		t.Fatalf("Expected empty pinned list on new storage, got: %v", loaded)
	}

	// Save items
	sample := []string{
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI user@phone",
		"Wi-Fi Password: SecretPassword123",
	}
	err = store.SavePinned(sample)
	if err != nil {
		t.Fatalf("Failed to save pinned clips: %v", err)
	}

	// Load items
	loaded = store.LoadPinned()
	if !slices.Equal(loaded, sample) {
		t.Fatalf("Loaded items mismatch: got %v, want %v", loaded, sample)
	}
}

func TestStorage_CorruptedFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "clipsync-corrupt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewStorageWithDir(tempDir)
	filePath := filepath.Join(tempDir, "pinned.json")

	// Write corrupted JSON
	_ = os.WriteFile(filePath, []byte("NOT_A_VALID_JSON{{{"), 0644)

	loaded := store.LoadPinned()
	if loaded != nil {
		t.Fatalf("Expected nil on corrupted file, got: %v", loaded)
	}
}
