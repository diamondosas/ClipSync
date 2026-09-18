package network

import (
	"clipsync/internal"
	"log"
	"time"

)

func PingIPS(ips []string){
	if len(ips) == 0 {
		return
	}

	for _, ip := range ips {
		sess, err :=  createPeer(ip)
		if err != nil {
			log.Println("SendPing Resolve Error:", err)
			return
		}

		msg := []byte{MsgTypePing}

		_, err = sess.Write(msg)
		if err != nil {
			log.Println("SendPing Write Error:", err)
			removePeer(ip)
		}

		log.Println("Sent Ping to", sess.RemoteAddr())
	}
}

func CheckForPing() {
	var activeIPS []string
	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if time.Since(internal.ConnDevices[i].LastSeen) <= (5 * time.Second) {
			internal.ConnDevices[i].Alive = true
			activeIPS= append(activeIPS, internal.ConnDevices[i].Ip)
		}else{
			log.Printf("[Ping] Device timed out: %s (%s)", internal.ConnDevices[i].Name, internal.ConnDevices[i].Ip)
			internal.ConnDevices[i].Alive = false
			removePeer(internal.ConnDevices[i].Ip)
		}
	}
}

