package service

import (
	"clipsync-android/internal"
	"clipsync-android/internal/events"
	"context"
	"testing"
	"time"
)

func TestClipSyncService_Lifecycle(t *testing.T) {
	bus := events.NewEventBus()
	svc := NewClipSyncService(bus, Config{
		Port:     0, // dynamic port for test
		Hostname: "TestPhoneService",
	})

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- svc.Start(ctx)
	}()

	// Let service start
	time.Sleep(50 * time.Millisecond)

	// Send manual broadcast
	svc.BroadcastClip("Test broadcast from service")

	// Connect manual IP (simulate)
	_ = svc.ConnectIP("127.0.0.1")

	// Cancel service
	cancel()

	select {
	case err := <-errCh:
		if err != nil && err != context.Canceled {
			t.Fatalf("Service stopped with unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Service did not shut down in time")
	}
}

func TestClipSyncService_PeerTracking(t *testing.T) {
	bus := events.NewEventBus()
	svc := NewClipSyncService(bus, Config{
		Port:     0,
		Hostname: "TestPhone",
	})

	svc.PeerManager().AddOrUpdatePeer(internal.Device{
		Name:     "DesktopPeer",
		IP:       "192.168.1.80",
		Alive:    true,
		LastSeen: time.Now(),
	})

	peers := svc.GetConnectedDevices()
	if len(peers) != 1 || peers[0].Name != "DesktopPeer" {
		t.Fatalf("Expected DesktopPeer in connected list, got: %+v", peers)
	}
}
