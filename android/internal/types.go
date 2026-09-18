package internal

import "time"

// Network and Protocol Constants
const (
	DefaultPort        = "9999"
	ServiceType        = "_clipsync._tcp"
	ServiceDomain      = "local."
	PingInterval       = 2 * time.Second
	PeerTimeout        = 5 * time.Second
	DiscoveryScanInterval = 3 * time.Second
)

// Global Secret Key for AES encryption across all ClipSync peers
var SecretKey = []byte("clipboardsyncapp")

// Message Protocol Types
const (
	MsgTypeHandshake byte = 0x01
	MsgTypePing      byte = 0x02
	MsgTypeClipboard byte = 0x03
)

// Device represents a remote network peer running ClipSync
type Device struct {
	Name     string
	IP       string
	Alive    bool
	LastSeen time.Time
}

// ClipItem represents a single clipboard history entry
type ClipItem struct {
	Content   string
	IsPinned  bool
	Timestamp time.Time
}
