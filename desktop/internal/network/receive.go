package network

import (
	"clipsync/internal"
	"clipsync/internal/view"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"log"
	"net"
	"time"
)

var (
	IncomingClips = make(chan []byte, 50)
)

const (
	MsgTypeHandshake byte = 0x01
	MsgTypePing      byte = 0x02
	MsgTypeClipboard byte = 0x03
)

// decryptCFB decrypts AES-CFB encrypted payloads where the first 16 bytes are the IV
func decryptCFB(data []byte, key []byte) ([]byte, error) {
	if len(data) < 16 {
		return nil, nil
	}
	iv := data[:16]
	ciphertext := data[16:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

// HandleIncomingPacket processes a datagram received on the UDP listener
func HandleIncomingPacket(data []byte, remoteAddr *net.UDPAddr) {
	if remoteAddr == nil || len(data) == 0 {
		return
	}

	ip := remoteAddr.IP.String()
	payload := data

	// If packet is encrypted (> 16 bytes and header not a raw message type), try decrypting
	if len(data) > 16 && data[0] != MsgTypeHandshake && data[0] != MsgTypePing && data[0] != MsgTypeClipboard {
		dec, err := decryptCFB(data, internal.SecretKey)
		if err == nil && len(dec) > 0 {
			payload = dec
		}
	}

	if len(payload) == 0 {
		return
	}

	msgType := payload[0]

	switch msgType {
	case MsgTypeHandshake:
		hostname := ""
		if len(payload) > 1 {
			hostname = string(payload[1:])
		}
		newDevice := internal.Device{
			Name:     hostname,
			Ip:       ip,
			LastSeen: time.Now(),
			Alive:    true,
		}
		view.AddNewDevice(newDevice)
		log.Printf("[Network] Received Handshake from %s (%s)", hostname, ip)

	case MsgTypePing:
		internal.ConnDevicesMu.Lock()
		found := false
		for i := range internal.ConnDevices {
			if ip == internal.ConnDevices[i].Ip {
				internal.ConnDevices[i].Alive = true
				internal.ConnDevices[i].LastSeen = time.Now()
				found = true
				break
			}
		}
		internal.ConnDevicesMu.Unlock()

		if !found {
			// Register unknown peer and exchange handshake
			view.AddNewDevice(internal.Device{
				Name:     ip,
				Ip:       ip,
				LastSeen: time.Now(),
				Alive:    true,
			})
			go Connect(ip)
		}
		log.Printf("[Network] %s sent Ping to Me", ip)

	case MsgTypeClipboard:
		content := payload[1:]
		internal.ConnDevicesMu.Lock()
		for i := range internal.ConnDevices {
			if ip == internal.ConnDevices[i].Ip {
				internal.ConnDevices[i].Alive = true
				internal.ConnDevices[i].LastSeen = time.Now()
			}
		}
		internal.ConnDevicesMu.Unlock()

		internal.LastRecvClipMu.Lock()
		internal.LastRecvClip = content
		internal.LastRecvClipMu.Unlock()

		// Push payload for the Clipboard goroutine to handle
		IncomingClips <- content

		log.Printf("[Network] Received clipboard (%d bytes) from %s", len(content), ip)

	default:
		log.Printf("[Network] Unknown message type 0x%x from %s", msgType, ip)
	}
}

// ReceiveClipboard returns the incoming clipboard channel
func ReceiveClipboard(ctx context.Context) <-chan []byte {
	return IncomingClips
}
