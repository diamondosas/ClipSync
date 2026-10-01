package network

import (
	"bytes"
	"clipsync/internal"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"net"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"
)

// Helper to encrypt plaintext with AES-CFB
func encryptCFB(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func TestDecryptCFB(t *testing.T) {
	// Change secretKey here to test custom AES keys (must be 16, 24, or 32 bytes)
	key := []byte("clipboardsyncapp")
	// Change message to test other decrypted payloads
	original := []byte("Sample sync payload")

	encrypted, err := encryptCFB(original, key)
	if err != nil {
		t.Fatalf("failed to encrypt test payload: %v", err)
	}

	decrypted, err := decryptCFB(encrypted, key)
	if err != nil {
		t.Fatalf("decryptCFB error: %v", err)
	}
	if !bytes.Equal(decrypted, original) {
		t.Errorf("decrypted %q; want %q", decrypted, original)
	}

	// Payloads under 16 bytes return nil without error
	short, err := decryptCFB([]byte("short"), key)
	if err != nil || short != nil {
		t.Errorf("expected nil for short payload, got %v, err=%v", short, err)
	}

	// Invalid key length returns error
	_, err = decryptCFB(encrypted, []byte("shortkey"))
	if err == nil {
		t.Error("expected error for invalid AES key length")
	}
}

func TestIsCleanText(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{"Clean ASCII", []byte("Hello ClipSync 123!"), true},
		{"With whitespace and tabs", []byte("Line 1\nLine 2\tEnd\r\n"), true},
		{"Unicode and emoji", []byte("ClipSync 🚀 剪贴板"), true},
		{"Empty byte slice", []byte(""), false},
		{"Binary control characters", []byte{0x00, 0x01, 0x02, 0x03}, false},
		{"Invalid UTF-8 sequence", []byte{0xff, 0xfe, 0xfd}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isCleanText(tc.input)
			if got != tc.expected {
				t.Errorf("isCleanText(%q) = %v; want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestHandleIncomingPacket_EdgeCases(t *testing.T) {
	// Calling with nil remoteAddr or empty data must return cleanly without panic
	HandleIncomingPacket(nil, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9999})
	HandleIncomingPacket([]byte("hello"), nil)
	HandleIncomingPacket([]byte{}, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9999})
}

func TestHandleIncomingPacket_Handshake(t *testing.T) {
	// Change peerName to test different hostnames
	peerName := "Office-Laptop"
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("192.168.1.101"), Port: 9999}
	packet := append([]byte{MsgTypeHandshake}, []byte(peerName)...)

	HandleIncomingPacket(packet, remoteAddr)

	internal.ConnDevicesMu.Lock()
	found := false
	for _, d := range internal.ConnDevices {
		if d.Ip == "192.168.1.101" && d.Name == peerName && d.Alive {
			found = true
			break
		}
	}
	internal.ConnDevicesMu.Unlock()

	if !found {
		t.Errorf("expected device %s (192.168.1.101) to be registered in ConnDevices", peerName)
	}
}

func TestHandleIncomingPacket_Ping(t *testing.T) {
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("192.168.1.102"), Port: 9999}
	packet := []byte{MsgTypePing}

	// 1. New unknown peer sending Ping should register
	HandleIncomingPacket(packet, remoteAddr)

	internal.ConnDevicesMu.Lock()
	var dev *internal.Device
	for i := range internal.ConnDevices {
		if internal.ConnDevices[i].Ip == "192.168.1.102" {
			dev = &internal.ConnDevices[i]
			break
		}
	}
	internal.ConnDevicesMu.Unlock()

	if dev == nil || !dev.Alive {
		t.Fatalf("expected device 192.168.1.102 to be registered and alive on Ping")
	}

	// 2. Existing peer sending Ping updates LastSeen
	prevSeen := dev.LastSeen
	time.Sleep(10 * time.Millisecond)
	HandleIncomingPacket(packet, remoteAddr)

	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if internal.ConnDevices[i].Ip == "192.168.1.102" {
			if !internal.ConnDevices[i].LastSeen.After(prevSeen) {
				t.Errorf("expected LastSeen to be updated on Ping")
			}
			break
		}
	}
	internal.ConnDevicesMu.Unlock()
}

func TestHandleIncomingPacket_Clipboard(t *testing.T) {
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("192.168.1.103"), Port: 9999}
	// Change clipText to test different incoming clipboard payloads
	clipText := "Incoming synchronized text"
	packet := append([]byte{MsgTypeClipboard}, []byte(clipText)...)

	HandleIncomingPacket(packet, remoteAddr)

	select {
	case received := <-ReceiveClipboard(context.Background()):
		if string(received) != clipText {
			t.Errorf("ReceiveClipboard() = %q; want %q", string(received), clipText)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for incoming clipboard packet")
	}

	// Test corrupt/binary clipboard payload gets discarded
	badPacket := append([]byte{MsgTypeClipboard}, []byte{0x00, 0x01, 0x02}...)
	HandleIncomingPacket(badPacket, remoteAddr)
}

func TestHandleIncomingPacket_Encrypted(t *testing.T) {
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("192.168.1.104"), Port: 9999}
	plaintext := append([]byte{MsgTypeClipboard}, []byte("Encrypted ClipSync Data")...)

	encrypted, err := encryptCFB(plaintext, internal.SecretKey)
	if err != nil {
		t.Fatalf("failed to encrypt packet: %v", err)
	}

	HandleIncomingPacket(encrypted, remoteAddr)

	select {
	case received := <-ReceiveClipboard(context.Background()):
		if string(received) != "Encrypted ClipSync Data" {
			t.Errorf("encrypted packet got %q; want %q", string(received), "Encrypted ClipSync Data")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for decrypted clipboard packet")
	}
}

func TestCheckForPing(t *testing.T) {
	internal.ConnDevicesMu.Lock()
	internal.ConnDevices = []internal.Device{
		{Name: "ActiveDev", Ip: "192.168.1.10", Alive: true, LastSeen: time.Now()},
		{Name: "DeadDev", Ip: "192.168.1.20", Alive: true, LastSeen: time.Now().Add(-15 * time.Second)},
	}
	internal.ConnDevicesMu.Unlock()

	CheckForPing()

	internal.ConnDevicesMu.Lock()
	defer internal.ConnDevicesMu.Unlock()
	for _, d := range internal.ConnDevices {
		if d.Ip == "192.168.1.10" && !d.Alive {
			t.Errorf("expected 192.168.1.10 to stay alive")
		}
		if d.Ip == "192.168.1.20" && d.Alive {
			t.Errorf("expected 192.168.1.20 to be marked inactive")
		}
	}
}

func TestPingIPS(t *testing.T) {
	// Pinging empty slice should return immediately
	PingIPS([]string{})

	// Start local mock UDP listener to receive ping
	listener, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start mock listener: %v", err)
	}
	defer listener.Close()

	_, portStr, _ := net.SplitHostPort(listener.LocalAddr().String())
	// Temporarily override PORT for test
	oldPort := internal.PORT
	internal.PORT = portStr
	defer func() { internal.PORT = oldPort }()

	PingIPS([]string{"127.0.0.1"})

	buf := make([]byte, 16)
	_ = listener.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	n, _, err := listener.ReadFrom(buf)
	if err != nil {
		t.Fatalf("mock listener failed to read ping: %v", err)
	}
	if n < 1 || buf[0] != MsgTypePing {
		t.Errorf("expected MsgTypePing (0x02), got 0x%x", buf[0])
	}
}

func TestSendClipboard(t *testing.T) {
	listener, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start mock listener: %v", err)
	}
	defer listener.Close()

	_, portStr, _ := net.SplitHostPort(listener.LocalAddr().String())
	oldPort := internal.PORT
	internal.PORT = portStr
	defer func() { internal.PORT = oldPort }()

	// Add mock target device
	internal.ConnDevicesMu.Lock()
	internal.ConnDevices = []internal.Device{
		{Name: "LocalReceiver", Ip: "127.0.0.1", Alive: true, LastSeen: time.Now()},
	}
	internal.ConnDevicesMu.Unlock()

	// Change testText to test sending other clipboard payloads
	testText := "Broadcast sync content"
	SendClipboard([]byte(testText))

	buf := make([]byte, 256)
	_ = listener.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	n, _, err := listener.ReadFrom(buf)
	if err != nil {
		t.Fatalf("failed to read clipboard packet: %v", err)
	}
	if n < 1 || buf[0] != MsgTypeClipboard {
		t.Fatalf("expected header MsgTypeClipboard (0x03), got 0x%x", buf[0])
	}
	if string(buf[1:n]) != testText {
		t.Errorf("got %q; want %q", string(buf[1:n]), testText)
	}
}

func TestConnect(t *testing.T) {
	listener, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start mock listener: %v", err)
	}
	defer listener.Close()

	_, portStr, _ := net.SplitHostPort(listener.LocalAddr().String())
	oldPort := internal.PORT
	internal.PORT = portStr
	defer func() { internal.PORT = oldPort }()

	internal.Hostname = "MyTestHost"
	Connect("127.0.0.1")

	buf := make([]byte, 256)
	_ = listener.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	n, _, err := listener.ReadFrom(buf)
	if err != nil {
		t.Fatalf("failed to read handshake packet: %v", err)
	}
	if buf[0] != MsgTypeHandshake {
		t.Errorf("expected MsgTypeHandshake (0x01), got 0x%x", buf[0])
	}
	if string(buf[1:n]) != "MyTestHost" {
		t.Errorf("got handshake host %q; want %q", string(buf[1:n]), "MyTestHost")
	}
}

func TestListen_Shutdown(t *testing.T) {
	// Change port temporarily to avoid conflict with standard 9999
	oldPort := internal.PORT
	internal.PORT = "19998"
	defer func() { internal.PORT = oldPort }()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := Listen(ctx)
	if err != nil {
		t.Errorf("Listen returned unexpected error: %v", err)
	}
}

func TestGetCurrIPS(t *testing.T) {
	// Verify GetCurrIPS executes without crashing
	ips := GetCurrIPS()
	t.Logf("Detected network signature IPs: %s", ips)
}

func TestEntry(t *testing.T) {
	results := make(chan *zeroconf.ServiceEntry, 1)
	entryDev := &zeroconf.ServiceEntry{
		Instance: "ZeroconfDevice",
		AddrIPv4: []net.IP{net.ParseIP("192.168.1.222")},
	}
	results <- entryDev
	close(results)

	entry(results)

	internal.ConnDevicesMu.Lock()
	defer internal.ConnDevicesMu.Unlock()
	found := false
	for _, d := range internal.ConnDevices {
		if d.Ip == "192.168.1.222" && d.Name == "ZeroconfDevice" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected zeroconf device to be added to ConnDevices")
	}
}
