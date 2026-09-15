package network

import (
	"clipsync/internal"
	"log"
	"slices"
	"sync"

	"github.com/xtaci/kcp-go/v5"
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
	if Sess == nil {
		log.Println("SendClipboard: Conn is nil, skipping.")
		return
	}

	internal.IPSMu.Lock()
	ips := make([]string, len(internal.IPS))
	copy(ips, internal.IPS)
	internal.IPSMu.Unlock()

	for _, ip := range ips {
		Sess , err:= kcp.DialWithOptions(ip + ":" + internal.PORT, BlockCrypt, 10, 3)
		if err != nil {
			log.Println("SendClipboard Resolve Error:", err)
			continue
		}
		_, err = Sess.Write(data)
		if err != nil {
			log.Println("SendClipboard Write Error:", err)
		}
	}
}
