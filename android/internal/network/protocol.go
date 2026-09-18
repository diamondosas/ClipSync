package network

import (
	"clipsync-android/internal"
	"errors"
	"fmt"
)

var (
	ErrEmptyPacket   = errors.New("empty packet received")
	ErrUnknownPacket = errors.New("unknown message type")
)

// EncodeHandshake wraps the device hostname in a MsgTypeHandshake packet.
func EncodeHandshake(hostname string) []byte {
	return append([]byte{internal.MsgTypeHandshake}, []byte(hostname)...)
}

// EncodePing creates a 1-byte MsgTypePing heartbeat packet.
func EncodePing() []byte {
	return []byte{internal.MsgTypePing}
}

// EncodeClipboard wraps clipboard text in a MsgTypeClipboard packet.
func EncodeClipboard(data []byte) []byte {
	return append([]byte{internal.MsgTypeClipboard}, data...)
}

// DecodePacket unpacks a raw wire packet into its message type byte and payload slice.
func DecodePacket(data []byte) (byte, []byte, error) {
	if len(data) == 0 {
		return 0, nil, ErrEmptyPacket
	}

	msgType := data[0]
	switch msgType {
	case internal.MsgTypeHandshake, internal.MsgTypePing, internal.MsgTypeClipboard:
		return msgType, data[1:], nil
	default:
		return 0, nil, fmt.Errorf("%w: 0x%x", ErrUnknownPacket, msgType)
	}
}
