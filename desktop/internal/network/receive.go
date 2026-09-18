package network

import (
	"clipsync/internal"
	"clipsync/internal/view"
	"context"
	"log"
	"net"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

var(  
	IncomingClips = make(chan []byte, 50)
)


const (
	MsgTypeHandshake byte = 0x01
	MsgTypePing      byte = 0x02
	MsgTypeClipboard byte = 0x03
)

func HandleIncomingData(sess *kcp.UDPSession) {
	defer sess.Close()

	ip := sess.RemoteAddr().(*net.UDPAddr).IP.String()
	
	tmpBuf := make([]byte, 65535)

	for{
		n, err := sess.Read(tmpBuf)

		if err != nil {
			log.Println("ReadFromUDP Error:", err)
			return 
		}

		MsgType := tmpBuf[0]
		
		payload := make([]byte, n)
		copy(payload, tmpBuf[1:n])


		switch MsgType{
		case MsgTypeHandshake:
			newDevice := internal.Device{
				Name: string(payload),
				Ip: ip,
				 LastSeen: time.Now(),
				Alive: true,
			}
			view.AddNewDevice(newDevice)

		case MsgTypePing:
			internal.ConnDevicesMu.Lock()
			for i := range internal.ConnDevices {
				if ip == string(internal.ConnDevices[i].Ip){
					internal.ConnDevices[i].Alive = true
					internal.ConnDevices[i].LastSeen = time.Now()
				}
			}
			internal.ConnDevicesMu.Unlock()
			log.Println(ip + "Sent Ping to Me")

		case MsgTypeClipboard:
			internal.LastRecvClipMu.Lock()
			internal.LastRecvClip = payload
			internal.LastRecvClipMu.Unlock()
			
			//Push the Payload for the Clipboard goroutine to handle
			IncomingClips <- payload

			log.Printf("[Network] Received clipboard (%d bytes) from %s", len(payload), ip)
		}
	} 
}

//Receive clipboard from peer and Update it
func ReceiveClipboard(ctx context.Context) chan []byte{
	for{
		select{
		case <-ctx.Done():
			return nil
		case <-IncomingClips:
			return IncomingClips
		}
	}
}



