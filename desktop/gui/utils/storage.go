// ai-generated
package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// getPinnedFilePath resolves the OS-specific path for pinned.json
// Linux: ~/.config/clipsync/pinned.json
// Windows: %APPDATA%\clipsync\pinned.json
// macOS: ~/Library/Application Support/clipsync/pinned.json
func getPinnedFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(configDir, "clipsync")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "pinned.json"), nil
}

// LoadPinned reads pinned clipboard items from disk.
func LoadPinned() []string {
	path, err := getPinnedFilePath()
	if err != nil {
		log.Printf("[Storage] Unable to resolve config dir: %v", err)
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		// File does not exist yet; normal on first launch
		return nil
	}
	var items []string
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("[Storage] Failed to parse pinned.json: %v", err)
		return nil
	}
	return items
}

// SavePinned saves pinned clipboard items to disk.
func SavePinned(items []string) error {
	path, err := getPinnedFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
