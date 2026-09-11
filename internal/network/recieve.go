package network

import (
	"clipsync/internal"
	"encoding/binary"
	"log"
	"net"
	"slices"
	"time"
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
	recievedData := tmpBuf[4 : 4+length]

	if slices.Equal(recievedData, []byte("---ClipSync---")) {
		UpdateIP()
	} else if slices.Equal(recievedData, []byte("~")) {
		UpdateDeviceState()
	} else {
		// Set LastRecievedClip to recievedData so other goroutines checking network.LastRecievedClip match correctly
		BufferMu.Lock()
		LastReceivedClip = make([]byte, len(recievedData))
		copy(LastReceivedClip, recievedData)
		BufferMu.Unlock()

		log.Println("Recieved Clipboard From Addr:", addr, "Content Length:", len(LastReceivedClip	))
		return recievedData, len(recievedData)
	}

	return nil, 0
}

func UpdateIP() {
	internal.IPSMu.Lock()
	found := false
	for _, existingIP := range internal.IPS {
		if existingIP == Addr.IP.String() {
			found = true
			break
		}
	}
	if !found {
		internal.IPS = append(internal.IPS, Addr.IP.String())
	}
	internal.IPSMu.Unlock()
}

func UpdateDeviceState() {
	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if string(Addr.IP.String()) == string(internal.ConnDevices[i].Ip){
			internal.ConnDevices[i].Alive = true
			internal.ConnDevices[i].LastSeen = time.Now()
		}
	}
	internal.ConnDevicesMu.Unlock()
	
}
