package network

import (
	"clipsync-android/internal"
	"clipsync-android/internal/events"
	"testing"
	"time"
)

func TestPeerManager_AddAndGetPeers(t *testing.T) {
	bus := events.NewEventBus()
	pm := NewPeerManager(bus)

	pm.AddOrUpdatePeer(internal.Device{
		Name:     "DesktopPC",
		IP:       "192.168.1.50",
		Alive:    true,
		LastSeen: time.Now(),
	})

	peers := pm.GetActivePeers()
	if len(peers) != 1 {
		t.Fatalf("Expected 1 active peer, got %d", len(peers))
	}
	if peers[0].Name != "DesktopPC" || peers[0].IP != "192.168.1.50" {
		t.Fatalf("Unexpected peer details: %+v", peers[0])
	}
}

func TestPeerManager_HeartbeatAndPrune(t *testing.T) {
	bus := events.NewEventBus()
	var disconnectedEvents []string
	sub := bus.SubscribeDevices(func(evt events.DeviceEvent) {
		if evt.Type == events.DeviceDisconnected {
			disconnectedEvents = append(disconnectedEvents, evt.Device.IP)
		}
	})
	defer sub.Unsubscribe()

	pm := NewPeerManager(bus)

	// Add an active peer
	pm.AddOrUpdatePeer(internal.Device{
		Name:     "ActivePeer",
		IP:       "192.168.1.10",
		Alive:    true,
		LastSeen: time.Now(),
	})

	// Add a timed-out peer (last seen 10s ago)
	pm.AddOrUpdatePeer(internal.Device{
		Name:     "OldPeer",
		IP:       "192.168.1.20",
		Alive:    true,
		LastSeen: time.Now().Add(-10 * time.Second),
	})

	// Run timeout pruning
	pruned := pm.PruneDeadPeers(5 * time.Second)
	if len(pruned) != 1 || pruned[0] != "192.168.1.20" {
		t.Fatalf("Expected pruned IP 192.168.1.20, got: %v", pruned)
	}

	active := pm.GetActivePeers()
	if len(active) != 1 || active[0].IP != "192.168.1.10" {
		t.Fatalf("Expected only active peer remaining, got: %+v", active)
	}

	time.Sleep(20 * time.Millisecond)
	if len(disconnectedEvents) != 1 || disconnectedEvents[0] != "192.168.1.20" {
		t.Fatalf("Expected disconnected event for 192.168.1.20, got: %v", disconnectedEvents)
	}
}

func TestPeerManager_TouchPeer(t *testing.T) {
	bus := events.NewEventBus()
	pm := NewPeerManager(bus)

	pm.AddOrUpdatePeer(internal.Device{
		Name:     "Peer1",
		IP:       "192.168.1.30",
		Alive:    true,
		LastSeen: time.Now().Add(-4 * time.Second),
	})

	// Touch peer on receiving ping
	pm.TouchPeer("192.168.1.30")

	peer, found := pm.GetPeer("192.168.1.30")
	if !found {
		t.Fatal("Expected peer to be found")
	}
	if time.Since(peer.LastSeen) > 1*time.Second {
		t.Fatalf("Expected LastSeen to be updated recently, was %v ago", time.Since(peer.LastSeen))
	}
}
