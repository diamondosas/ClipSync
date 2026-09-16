package network

import (
	"clipsync/internal"
	"log"
	"time"


	"github.com/xtaci/kcp-go/v5"
)

func PingIPS(ips []string) {
	if len(ips) == 0 {
		return
	}

	for _, ip := range ips {
		SendPing(ip)
	}
}

func SendPing(ip string) {
	sess, err :=  kcp.DialWithOptions(ip + ":" + internal.PORT, BlockCrypt, 10, 3)
	if err != nil {
		log.Println("SendPing Resolve Error:", err)
		return
	}

	msg := []byte{MsgTypePing}

	_, err = sess.Write(msg)
	if err != nil {
		log.Println("SendPing Write Error:", err)
	}

	log.Println("Sent Ping to", sess.RemoteAddr())
}

func CheckForPing() {
	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if time.Since(internal.ConnDevices[i].LastSeen) >= (time.Second * 5) {
			log.Println("Dead connection Found", )
			//Remove it from internal completely and make sure to send the signal to the GUI
			internal.ConnDevices[i].Alive = false
			removePeer(internal.ConnDevices[i].Ip)
		}
	}

	internal.ConnDevicesMu.Unlock()
}
