package network

import (
	"clipsync/internal"
	"clipsync/internal/view"
	"context"
	"log"
	"net"
	"slices"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

var(  
	IncomingClips = make(chan []byte, 50)
)

const (
	MsgTypeHandshake byte = 0x01
	MsgTypePing      byte = 0x02
)

func HandleIncomingData(sess *kcp.UDPSession) {
	defer sess.Close()

	ip := getIP(sess)
	
	tmpBuf := make([]byte, 65535)

	for{
		n, err := sess.Read(tmpBuf)

		if err != nil {
			log.Println("ReadFromUDP Error:", err)
			return 
		}

		receivedData := make([]byte, n)
		copy(receivedData, tmpBuf[:n])

		if len(receivedData) > 0 && receivedData[0] == MsgTypeHandshake{
			UpdateIP(ip)
			view.AddNewDevice(
				internal.Device{
				Name: string(receivedData[1:]),
				Ip: ip,
				LastSeen: time.Now(),
				Alive: true,
			})
			go Connect(ip)
		} else if slices.Equal(receivedData, []byte{MsgTypePing}) {
			UpdateDeviceState(ip)
		} else {
			// Set LastRecievedClip to receivedData so other goroutines checking network.LastRecievedClip match correctly
			BufferMu.Lock()
			LastReceivedClip = make([]byte, len(receivedData))
			copy(LastReceivedClip, receivedData)
			BufferMu.Unlock()

			log.Printf("[Network] Received clipboard (%d bytes) from %s", len(receivedData), ip)
			IncomingClips <- receivedData
		}

	} 
}

func ReceiveClipboard(ctx context.Context) []byte{
	for{
		select{
		case <-ctx.Done():
			return nil
		case clip :=  <-IncomingClips:
			return clip
		}
	}
}


func getIP(sess *kcp.UDPSession) string{
	remoteAddr, ok := sess.RemoteAddr().(*net.UDPAddr)
	if !ok || remoteAddr == nil{
		return ""
	}

	return remoteAddr.IP.String()
}

func UpdateIP(ip string) {
	internal.IPSMu.Lock()
	found := false
	for _, existingIP := range internal.IPS {
		if existingIP == ip {
			found = true
			break
		}
	}
	if !found {
		internal.IPS = append(internal.IPS, ip)
	}
	internal.IPSMu.Unlock()
	log.Println("Acknoledge Handshake")
}

func UpdateDeviceState(ip string) {
	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if ip == string(internal.ConnDevices[i].Ip){
			internal.ConnDevices[i].Alive = true
			internal.ConnDevices[i].LastSeen = time.Now()
		}
	}
	internal.ConnDevicesMu.Unlock()
	log.Println(ip + "Sent Ping to Me")
}
