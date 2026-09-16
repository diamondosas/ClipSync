package network

import (
	"clipsync/internal"
	"log"
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
	internal.IPSMu.Lock()
	ips := make([]string, len(internal.IPS))
	copy(ips, internal.IPS)
	internal.IPSMu.Unlock()

	for _, ip := range ips {
		sess , err := createPeer(ip)
		if err != nil || sess == nil{
			continue
		}
		_, err = sess.Write(data)
		if err != nil {
			log.Println("SendClipboard Write Error:", err)
			removePeer(ip)
			return
		}

	}
}
