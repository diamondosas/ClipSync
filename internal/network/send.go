package network

import (
	"clipsync/internal"
	"encoding/binary"
	"log"
	"net"
	"slices"
	"sync"
)

var (
	BufferMu         sync.RWMutex
	LastReceivedClip []byte
)

// IsLastReceived checks whether data matches the last received clipboard buffer in a thread-safe way.
func IsLastReceived(data []byte) bool {
	BufferMu.RLock()
	defer BufferMu.RUnlock()
	return slices.Equal(data, LastReceivedClip)
}

func SendClipboard(data []byte) {
	if Conn == nil {
		log.Println("SendClipboard: Conn is nil, skipping.")
		return
	}

	payload := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(payload[:4], uint32(len(data)))
	copy(payload[4:], data)

	internal.IPSMu.Lock()
	ips := make([]string, len(internal.IPS))
	copy(ips, internal.IPS)
	internal.IPSMu.Unlock()

	for _, ip := range ips {
		addr, err := net.ResolveUDPAddr("udp", ip+":"+internal.PORT)
		if err != nil {
			log.Println("SendClipboard Resolve Error:", err)
			continue
		}
		_, err = Conn.WriteToUDP(payload, addr)
		if err != nil {
			log.Println("SendClipboard Write Error:", err)
		}
	}
}
