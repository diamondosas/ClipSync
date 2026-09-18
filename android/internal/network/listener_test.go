package network

import (
	"clipsync-android/internal"
	"clipsync-android/internal/events"
	"context"
	"net"
	"testing"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

func TestListener_HandleIncomingPackets(t *testing.T) {
	bus := events.NewEventBus()
	pm := NewPeerManager(bus)

	var receivedClips []string
	var connectedDevices []string

	sub1 := bus.SubscribeClips(func(evt events.ClipReceivedEvent) {
		receivedClips = append(receivedClips, evt.Content)
	})
	defer sub1.Unsubscribe()

	sub2 := bus.SubscribeDevices(func(evt events.DeviceEvent) {
		if evt.Type == events.DeviceConnected {
			connectedDevices = append(connectedDevices, evt.Device.Name)
		}
	})
	defer sub2.Unsubscribe()

	listener := NewListener(pm, bus, internal.SecretKey, 0) // dynamic port 0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	boundPort, err := listener.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}

	block, _ := kcp.NewAESBlockCrypt(internal.SecretKey)
	sess, err := kcp.DialWithOptions(net.JoinHostPort("127.0.0.1", boundPort), block, 10, 3)
	if err != nil {
		t.Fatalf("Failed to dial test listener: %v", err)
	}
	defer sess.Close()

	// 1. Send Handshake
	_, err = sess.Write(EncodeHandshake("RemoteTestHost"))
	if err != nil {
		t.Fatalf("Failed to write handshake: %v", err)
	}

	// 2. Send Ping
	_, err = sess.Write(EncodePing())
	if err != nil {
		t.Fatalf("Failed to write ping: %v", err)
	}

	// 3. Send Clipboard
	_, err = sess.Write(EncodeClipboard([]byte("Network Clipboard Data")))
	if err != nil {
		t.Fatalf("Failed to write clipboard: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if len(connectedDevices) == 0 || connectedDevices[0] != "RemoteTestHost" {
		t.Fatalf("Expected connected device 'RemoteTestHost', got: %v", connectedDevices)
	}

	if len(receivedClips) == 0 || receivedClips[0] != "Network Clipboard Data" {
		t.Fatalf("Expected received clip 'Network Clipboard Data', got: %v", receivedClips)
	}
}
