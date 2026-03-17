package network_test

import (
	"clipsync/internal/globals"
	"clipsync/internal/network"
	"context"
	"os"
	"testing"
	"time"
)

func TestRegisterAndBrowse(t *testing.T) {
	// Set our browser's name so we don't filter out the test registration
	globals.Username, _ = os.Hostname()
	testServiceName := "ClipSync-Test-Device"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Start registering a service with a unique name
	go func() {
		if err := network.RegisterDevice(ctx, testServiceName); err != nil {
			t.Errorf("RegisterDevice failed: %v", err)
		}
	}()

	// Give zeroconf some time to start broadcasting
	time.Sleep(2 * time.Second)

	// 2. Clear existing IPs and start browsing
	globals.IP = nil
	if err := network.BrowseForDevices(); err != nil {
		t.Fatalf("BrowseForDevices failed: %v", err)
	}

	// 3. Wait for discovery (max 10 seconds)
	found := false
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for !found {
		select {
		case <-timeout:
			t.Fatal("Timed out waiting for device discovery")
		case <-ticker.C:
			if len(globals.IP) > 0 {
				found = true
				t.Logf("Found devices: %v", globals.IP)
			}
		}
	}

	if !found {
		t.Error("Did not find any devices during browse")
	}
}
