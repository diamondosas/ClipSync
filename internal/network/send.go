package network

import (
	"encoding/binary"
	"log"
	"net"
	"slices"
	"sync"

	"clipsync/internal/globals"
)

var (
	BufferMu sync.RWMutex
	LastRecievedClip   []byte
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
	
	globals.IPSMu.Lock()
	ips := make([]string, len(globals.IPS))
	copy(ips, globals.IPS)
	globals.IPSMu.Unlock()

	for _, ip := range ips {
		addr, err := net.ResolveUDPAddr("udp", ip + ":" + globals.PORT)
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

