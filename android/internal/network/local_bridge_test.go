package network

import (
	"clipsync-android/internal/events"
	"context"
	"net"
	"testing"
	"time"
)

func TestLocalBridge_ReceivesAccessibilityClipAndBroadcasts(t *testing.T) {
	bus := events.NewEventBus()
	pm := NewPeerManager(bus)

	var receivedClips []string
	sub := bus.SubscribeClips(func(evt events.ClipReceivedEvent) {
		receivedClips = append(receivedClips, evt.Content)
	})
	defer sub.Unsubscribe()

	bridge := NewLocalBridge(bus, pm, nil, 0) // dynamic port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	boundPort, err := bridge.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start local bridge: %v", err)
	}

	// Simulate Android Accessibility Service sending copied text to localhost
	conn, err := net.Dial("udp", net.JoinHostPort("127.0.0.1", boundPort))
	if err != nil {
		t.Fatalf("Failed to dial local bridge: %v", err)
	}
	defer conn.Close()

	copiedText := "Copied password or URL from Android browser"
	_, err = conn.Write([]byte(copiedText))
	if err != nil {
		t.Fatalf("Failed to send local clip: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if len(receivedClips) == 0 || receivedClips[0] != copiedText {
		t.Fatalf("Expected local bridge to dispatch clip %q, got: %v", copiedText, receivedClips)
	}
}
