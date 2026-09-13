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

const (
	MsgTypeHandshake byte = 0x01
	MsgTypePing      byte = 0x02
)

func ReceiveData() ([]byte, int) {
	if Conn == nil {
		log.Println("ReceiveClipboard: Conn is nil. Waiting for Ready...")
		<-Ready
	}

	tmpBuf := make([]byte, 65535)
	n, addr, err := Conn.ReadFromUDP(tmpBuf)
	Addr = addr

	if err != nil || n < 4 {
		log.Println("ReadFromUDP Error:", err)
		return nil, 0
	}

	length := binary.BigEndian.Uint32(tmpBuf[:4])
	receivedData := tmpBuf[4 : 4+length]

	if slices.Equal(receivedData, []byte{MsgTypeHandshake}){
		UpdateIP()
	} else if slices.Equal(receivedData, []byte{MsgTypePing}) {
		UpdateDeviceState()
	} else {
		// Set LastRecievedClip to receivedData so other goroutines checking network.LastRecievedClip match correctly
		BufferMu.Lock()
		LastReceivedClip = make([]byte, len(receivedData))
		copy(LastReceivedClip, receivedData)
		BufferMu.Unlock()

		log.Println("Recieved Clipboard From Addr:", addr, "Content Length:", len(LastReceivedClip))
		return receivedData, len(receivedData)
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
		if string(Addr.IP.String()) == string(internal.ConnDevices[i].Ip) {
			internal.ConnDevices[i].Alive = true
			internal.ConnDevices[i].LastSeen = time.Now()
		}
	}
	internal.ConnDevicesMu.Unlock()

}
