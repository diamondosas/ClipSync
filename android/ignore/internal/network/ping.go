package network

import (
	"clipsync/internal"
	"clipsync/internal/view"
	"log"
	"time"
)

func PingIPS(ips []string) {
	if len(ips) == 0 {
		return
	}

	msg := []byte{MsgTypePing}

	for _, ip := range ips {
		err := SendPacket(ip, msg)
		if err != nil {
			log.Println("SendPing Write Error to", ip, ":", err)
			continue
		}

		log.Println("[Ping] Sent Ping to", ip)
	}
}

func CheckForPing() {
	var activeIPS []string
	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if time.Since(internal.ConnDevices[i].LastSeen) <= (10 * time.Second) {
			internal.ConnDevices[i].Alive = true
			activeIPS = append(activeIPS, internal.ConnDevices[i].Ip)
		} else {
			log.Printf("[Ping] Device timed out: %s (%s)", internal.ConnDevices[i].Name, internal.ConnDevices[i].Ip)
			internal.ConnDevices[i].Alive = false
		}
	}
	internal.ConnDevicesMu.Unlock()

	view.UpdateDevices(activeIPS)
}
