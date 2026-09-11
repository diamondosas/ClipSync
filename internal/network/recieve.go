package network

import(
	"encoding/binary"
	"log"
	"slices"
	"net"

	"clipsync/internal/globals"
)

var Addr *net.UDPAddr

func RecieveData() ([]byte, int) {
	if Conn == nil {
		log.Println("RecieveClipboard: Conn is nil. Waiting for Ready...")
		<-Ready
	}
	
	tmpBuf := make([]byte, 65535)
	_, addr, err := Conn.ReadFromUDP(tmpBuf)
	Addr = addr

	if err != nil {
		log.Println("ReadFromUDP Error:", err)
		return nil, 0
	}

	length := binary.BigEndian.Uint32(tmpBuf[:4])
	recievedData := tmpBuf[4 : 4 + length]
	


	if slices.Equal(recievedData, []byte("---ClipSync---")) {
		UpdateIP()
	} else if slices.Equal(recievedData, []byte("~")) {
		UpdateDeviceState()
	} else {
		// Set LastRecievedClip to recievedData so other goroutines checking network.LastRecievedClip match correctly
		BufferMu.Lock()
		LastRecievedClip = make([]byte, len(recievedData))
		copy(LastRecievedClip, recievedData)
		BufferMu.Unlock()

		log.Println("Recieved Clipboard From Addr:", addr, "Content Length:", len(LastRecievedClip))
		return recievedData, len(recievedData)
	}

	return nil, 0
}


func UpdateIP(){
	globals.IPSMu.Lock()
	found := false
	for _, existingIP := range globals.IPS {
		if existingIP == Addr.IP.String() {
			found = true
			break
		}
	}
	if !found {
		globals.IPS = append(globals.IPS, Addr.IP.String())
	}
	globals.IPSMu.Unlock()
}

func UpdateDeviceState(){
	globals.ConnDevicesMu.Lock()
	for i := range globals.ConnDevices {
		globals.ConnDevices[i].Alive = true
	}
	globals.ConnDevicesMu.Unlock()
}