package network

import (
	"bytes"
	"clipsync-android/internal"
	"testing"
)

func TestEncodeHandshake(t *testing.T) {
	hostname := "AndroidDevice-S23"
	packet := EncodeHandshake(hostname)

	if len(packet) == 0 {
		t.Fatal("Expected non-empty packet")
	}
	if packet[0] != internal.MsgTypeHandshake {
		t.Fatalf("Expected MsgTypeHandshake (0x01), got: 0x%x", packet[0])
	}
	if string(packet[1:]) != hostname {
		t.Fatalf("Expected hostname %q, got: %q", hostname, string(packet[1:]))
	}
}

func TestEncodePing(t *testing.T) {
	packet := EncodePing()
	if len(packet) != 1 {
		t.Fatalf("Expected 1 byte ping packet, got: %d bytes", len(packet))
	}
	if packet[0] != internal.MsgTypePing {
		t.Fatalf("Expected MsgTypePing (0x02), got: 0x%x", packet[0])
	}
}

func TestEncodeClipboard(t *testing.T) {
	clipData := "https://example.com/share/token123"
	packet := EncodeClipboard([]byte(clipData))

	if len(packet) != 1+len(clipData) {
		t.Fatalf("Expected packet length %d, got: %d", 1+len(clipData), len(packet))
	}
	if packet[0] != internal.MsgTypeClipboard {
		t.Fatalf("Expected MsgTypeClipboard (0x03), got: 0x%x", packet[0])
	}
	if string(packet[1:]) != clipData {
		t.Fatalf("Expected payload %q, got: %q", clipData, string(packet[1:]))
	}
}

func TestDecodePacket(t *testing.T) {
	// 1. Handshake
	hsPacket := []byte{internal.MsgTypeHandshake, 'L', 'i', 'n', 'u', 'x', 'P', 'C'}
	msgType, payload, err := DecodePacket(hsPacket)
	if err != nil {
		t.Fatalf("Unexpected error decoding handshake: %v", err)
	}
	if msgType != internal.MsgTypeHandshake || string(payload) != "LinuxPC" {
		t.Fatalf("Handshake decode mismatch: type=%d payload=%q", msgType, string(payload))
	}

	// 2. Ping
	pingPacket := []byte{internal.MsgTypePing}
	msgType, payload, err = DecodePacket(pingPacket)
	if err != nil {
		t.Fatalf("Unexpected error decoding ping: %v", err)
	}
	if msgType != internal.MsgTypePing || len(payload) != 0 {
		t.Fatalf("Ping decode mismatch: type=%d len=%d", msgType, len(payload))
	}

	// 3. Clipboard
	clipPacket := append([]byte{internal.MsgTypeClipboard}, []byte("Test Clip Data")...)
	msgType, payload, err = DecodePacket(clipPacket)
	if err != nil {
		t.Fatalf("Unexpected error decoding clipboard: %v", err)
	}
	if msgType != internal.MsgTypeClipboard || !bytes.Equal(payload, []byte("Test Clip Data")) {
		t.Fatalf("Clipboard decode mismatch: type=%d payload=%q", msgType, string(payload))
	}

	// 4. Empty packet
	_, _, err = DecodePacket([]byte{})
	if err == nil {
		t.Fatal("Expected error on empty packet, got nil")
	}

	// 5. Unknown packet type
	_, _, err = DecodePacket([]byte{0xFF, 't', 'e', 's', 't'})
	if err == nil {
		t.Fatal("Expected error on unknown packet type, got nil")
	}
}
