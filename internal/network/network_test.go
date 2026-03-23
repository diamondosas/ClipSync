package network_test

import (
	"clipsync/internal/globals"
	"clipsync/internal/network"
	"testing"
	"time"
)

func TestFullNetworkWorkflow(t *testing.T) {
	// 1. Setup context and global variables
	ctx := t.Context()

	// Set our browser's name so we don't filter out the test registration
	// We want to make sure it's different from the test device name
	globals.Username = "Test-Browser"
	testDeviceName := "ClipSync-Test-Device"
	globals.IPS = nil

	// 2. Start Listening
	go func() {
		if err := network.Listen(ctx); err != nil {
			t.Logf("Listen stopped: %v", err)
		}
	}()

	// Wait for listener to be ready
	select {
	case <-network.Ready:
		t.Log("Listener is ready")
	case <-time.After(5 * time.Second):
		t.Fatal("Listener timed out waiting to be ready")
	}

	// 3. Start registering a service with a unique name
	go func() {
		if err := network.RegisterDevice(ctx, testDeviceName); err != nil {
			t.Errorf("RegisterDevice failed: %v", err)
		}
	}()

	// Give zeroconf some time to start broadcasting
	time.Sleep(2 * time.Second)

	// 4. Start browsing for devices
	go func() {
		if err := network.BrowseForDevices(ctx); err != nil {
			t.Logf("BrowseForDevices stopped: %v", err)
		}
	}()

	// 5. Wait for discovery and connection verification
	// We expect:
	// - globals.IPS to be updated
	// - A "ClipSync" message to be received via UDP because BrowseForDevices calls Connect(ip)
	
	found := false
	receivedCS := false
	timeout := time.After(15 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	received := make(chan string)
	go func() {
		for {
			buf, n := network.RecieveClipboard()
			if n > 0 {
				received <- string(buf[:n])
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-timeout:
			if !found {
				t.Error("Timed out waiting for device discovery")
			}
			if !receivedCS {
				t.Error("Timed out waiting for ClipSync handshake message")
			}
			return
		case msg := <-received:
			t.Logf("Received message: %s", msg)
			if msg == "ClipSync" {
				t.Log("Successfully received ClipSync handshake")
				receivedCS = true
			}
			if found && receivedCS {
				t.Log("All conditions met: device found and handshake received.")
				return
			}
		case <-ticker.C:
			if len(globals.IPS) > 0 {
				if !found {
					found = true
					t.Logf("Found devices: %v", globals.IPS)
				}
			}
			if found && receivedCS {
				t.Log("All conditions met: device found and handshake received.")
				return
			}
		}
	}
}
