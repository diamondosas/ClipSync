package network

import (
	"encoding/binary"
	"log"
	"net"
	"slices"
	"sync"
)

var (
	BufferMu         sync.RWMutex
	LastRecievedClip []byte
)

// IsLastReceived checks whether data matches the last received clipboard buffer in a thread-safe way.
func IsLastReceived(data []byte) bool {
	BufferMu.RLock()
	defer BufferMu.RUnlock()
	return slices.Equal(data, LastRecievedClip)
}

func SendClipboard(data []byte) {
	if Conn == nil {
		log.Println("SendClipboard: Conn is nil, skipping.")
		return
	}

	payload := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(payload[:4], uint32(len(data)))
	copy(payload[4:], data)

	IPSMu.Lock()
	ips := make([]string, len(IPS))
	copy(ips, IPS)
	IPSMu.Unlock()

	for _, ip := range ips {
		addr, err := net.ResolveUDPAddr("udp", ip+":"+PORT)
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
