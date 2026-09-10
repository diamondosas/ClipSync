package network

import (
	"encoding/binary"
	"log"
	"net"
	"slices"
	"strconv"
	"sync"

	"clipsync/internal/globals"
)

var (
	BufferMu sync.RWMutex
	Buffer   []byte
)

// IsLastReceived checks whether data matches the last received clipboard buffer in a thread-safe way.
func IsLastReceived(data []byte) bool {
	BufferMu.RLock()
	defer BufferMu.RUnlock()
	return slices.Equal(data, Buffer)
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
		addr, err := net.ResolveUDPAddr("udp", ip + ":" + strconv.Itoa(globals.PORT))
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

func RecieveClipboard() ([]byte, int) {
	if Conn == nil {
		log.Println("RecieveClipboard: Conn is nil. Waiting for Ready...")
		<-Ready
	}
	tmpBuf := make([]byte, 65535)
	n, addr, err := Conn.ReadFromUDP(tmpBuf)
	if err != nil {
		log.Println("ReadFromUDP Error:", err)
		return nil, 0
	}
	
	if n < 4 {
		return nil, 0
	}
	
	length := binary.BigEndian.Uint32(tmpBuf[:4])
	if length > uint32(n-4) {
		log.Println("Incomplete payload received")
		return nil, 0
	}
	
	actualData := tmpBuf[4 : 4+length]
	
	if slices.Equal(actualData, []byte("---ClipSync---")) {
		globals.IPSMu.Lock()
		found := false
		for _, existingIP := range globals.IPS {
			if existingIP == addr.IP.String() {
				found = true
				break
			}
		}
		if !found {
			globals.IPS = append(globals.IPS, addr.IP.String())
		}
		globals.IPSMu.Unlock()
	} else if slices.Equal(actualData, []byte("---Ping---")) {
		SendPong(addr)
	} else if slices.Equal(actualData, []byte("---Pong---")) {
		select {
		case PongChan <- addr.IP.String():
		default:
		}
	} else {
		// Set Buffer to actualData so other goroutines checking network.Buffer match correctly
		BufferMu.Lock()
		Buffer = make([]byte, len(actualData))
		copy(Buffer, actualData)
		BufferMu.Unlock()
		log.Println("Recieved Clipboard From Addr:", addr, "Content Length:", len(Buffer))
		return actualData, len(actualData)
	}

	return nil, 0
}
