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
	if Sess == nil {
		return
	}

	Sess, err :=  kcp.DialWithOptions(ip + ":" + internal.PORT, BlockCrypt, 10, 3)
if err != nil {
		log.Println("SendPing Resolve Error:", err)
		return
	}

	msg := []byte{MsgTypePing}

	_, err = Sess.Write(msg)
	if err != nil {
		log.Println("SendPing Write Error:", err)
	}
}

func CheckForPing() {
	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if time.Since(internal.ConnDevices[i].LastSeen) >= (time.Second * 5) {
			internal.ConnDevices[i].Alive = false
		}
	}
	internal.ConnDevicesMu.Unlock()
}
