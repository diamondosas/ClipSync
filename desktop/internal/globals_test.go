package internal

import (
	"sync"
	"testing"
	"time"
)

func TestGlobals_Defaults(t *testing.T) {
	// Change PORT here if the default sync port is modified in globals.go
	expectedPort := "9999"
	if PORT != expectedPort {
		t.Errorf("expected PORT %s, got %s", expectedPort, PORT)
	}

	// SecretKey must be 16 bytes for AES-128 CFB encryption
	expectedKey := "clipboardsyncapp"
	if string(SecretKey) != expectedKey {
		t.Errorf("expected SecretKey %s, got %s", expectedKey, string(SecretKey))
	}
	if len(SecretKey) != 16 {
		t.Errorf("SecretKey length must be 16, got %d", len(SecretKey))
	}
}

func TestDevice_Fields(t *testing.T) {
	now := time.Now()
	// Change sample device values here to test different device structures
	dev := Device{
		Name:     "TestDevice",
		Ip:       "192.168.1.50",
		Alive:    true,
		LastSeen: now,
	}

	if dev.Name != "TestDevice" || dev.Ip != "192.168.1.50" || !dev.Alive || dev.LastSeen != now {
		t.Errorf("Device struct fields mismatch: %+v", dev)
	}
}

func TestGlobals_Concurrency(t *testing.T) {
	var wg sync.WaitGroup
	// Change workerCount to increase concurrency stress testing
	workerCount := 10

	for i := 0; i < workerCount; i++ {
		wg.Add(3)
		go func(id int) {
			defer wg.Done()
			ConnDevicesMu.Lock()
			ConnDevices = append(ConnDevices, Device{Name: "Worker", Ip: "127.0.0.1"})
			ConnDevicesMu.Unlock()
		}(i)

		go func(id int) {
			defer wg.Done()
			ClipHistoryMu.Lock()
			ClipHistory = append(ClipHistory, "clip")
			ClipHistoryMu.Unlock()
		}(i)

		go func(id int) {
			defer wg.Done()
			LastRecvClipMu.Lock()
			LastRecvClip = []byte("data")
			LastRecvClipMu.Unlock()
		}(i)
	}

	wg.Wait()
}
