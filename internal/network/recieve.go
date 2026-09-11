package network

import(
	"encoding/binary"
	"log"
	"slices"

	"clipsync/internal/globals"
)

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
		// Set LastRecievedClip to actualData so other goroutines checking network.LastRecievedClip match correctly
		BufferMu.Lock()
		LastRecievedClip = make([]byte, len(actualData))
		copy(LastRecievedClip, actualData)
		BufferMu.Unlock()
		log.Println("Recieved Clipboard From Addr:", addr, "Content Length:", len(LastRecievedClip))
		return actualData, len(actualData)
	}

	return nil, 0
}
