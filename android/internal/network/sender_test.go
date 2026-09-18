package network

import (
	"clipsync-android/internal"
	"clipsync-android/internal/events"
	"net"
	"testing"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

func TestSender_SendAndReceiveIntegration(t *testing.T) {
	// Setup a real local KCP listener on a random available port for testing
	block, err := kcp.NewAESBlockCrypt(internal.SecretKey)
	if err != nil {
		t.Fatalf("Failed to create crypt: %v", err)
	}

	listener, err := kcp.ListenWithOptions("127.0.0.1:0", block, 10, 3)
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.UDPAddr).Port

	bus := events.NewEventBus()
	pm := NewPeerManager(bus)

	sender := NewSender(pm, internal.SecretKey, port)
	defer sender.CloseAll()

	// Server accept goroutine
	receivedPayload := make(chan []byte, 1)
	go func() {
		sess, err := listener.AcceptKCP()
		if err != nil {
			return
		}
		defer sess.Close()
		buf := make([]byte, 1024)
		n, err := sess.Read(buf)
		if err == nil && n > 0 {
			receivedPayload <- buf[:n]
		}
	}()

	// Connect to our test listener
	err = sender.SendHandshake("127.0.0.1", "AndroidTestClient")
	if err != nil {
		t.Fatalf("SendHandshake failed: %v", err)
	}

	select {
	case data := <-receivedPayload:
		if len(data) == 0 || data[0] != internal.MsgTypeHandshake || string(data[1:]) != "AndroidTestClient" {
			t.Fatalf("Unexpected handshake packet received: %v", data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for handshake packet")
	}
}

func TestSender_BroadcastClipboard(t *testing.T) {
	block, err := kcp.NewAESBlockCrypt(internal.SecretKey)
	if err != nil {
		t.Fatalf("Failed to create crypt: %v", err)
	}

	listener, err := kcp.ListenWithOptions("127.0.0.1:0", block, 10, 3)
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.UDPAddr).Port

	bus := events.NewEventBus()
	pm := NewPeerManager(bus)
	pm.AddOrUpdatePeer(internal.Device{
		Name:     "TestDevice",
		IP:       "127.0.0.1",
		Alive:    true,
		LastSeen: time.Now(),
	})

	sender := NewSender(pm, internal.SecretKey, port)
	defer sender.CloseAll()

	receivedClip := make(chan string, 1)
	go func() {
		sess, err := listener.AcceptKCP()
		if err != nil {
			return
		}
		defer sess.Close()
		buf := make([]byte, 1024)
		n, err := sess.Read(buf)
		if err == nil && n > 0 {
			if buf[0] == internal.MsgTypeClipboard {
				receivedClip <- string(buf[1:n])
			}
		}
	}()

	sender.BroadcastClipboard([]byte("Copied text from test"))

	select {
	case clip := <-receivedClip:
		if clip != "Copied text from test" {
			t.Fatalf("Expected 'Copied text from test', got %q", clip)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for clipboard packet")
	}
}
