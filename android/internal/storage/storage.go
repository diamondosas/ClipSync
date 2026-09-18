package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Storage struct {
	customDir string
}

func NewStorage() *Storage {
	return &Storage{}
}

func NewStorageWithDir(dir string) *Storage {
	return &Storage{customDir: dir}
}

func (s *Storage) getFilePath() (string, error) {
	if s.customDir != "" {
		if err := os.MkdirAll(s.customDir, 0755); err != nil {
			return "", err
		}
		return filepath.Join(s.customDir, "pinned.json"), nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" || configDir == "." {
		if home := os.Getenv("HOME"); home != "" {
			configDir = filepath.Join(home, "config")
		} else {
			configDir = filepath.Join(os.TempDir(), "clipsync-config")
		}
	}
	dir := filepath.Join(configDir, "clipsync")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "pinned.json"), nil
}

func (s *Storage) LoadPinned() []string {
	path, err := s.getFilePath()
	if err != nil {
		log.Printf("[Storage] Unable to resolve storage path: %v", err)
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var items []string
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("[Storage] Failed to parse pinned.json: %v", err)
		return nil
	}
	return items
}

func (s *Storage) SavePinned(items []string) error {
	path, err := s.getFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
