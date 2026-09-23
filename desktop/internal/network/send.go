package network

import (
	"clipsync/internal"
	"log"
)

// SendClipboard sends clipboard data to all known alive connected peers
func SendClipboard(data []byte) {
	var ips []string

	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if internal.ConnDevices[i].Alive {
			ips = append(ips, internal.ConnDevices[i].Ip)
		}
	}
	internal.ConnDevicesMu.Unlock()

	payload := append([]byte{MsgTypeClipboard}, data...)

	for _, ip := range ips {
		err := SendPacket(ip, payload)
		if err != nil {
			log.Println("SendClipboard Write Error to", ip, ":", err)
			continue
		}
		log.Printf("[Network] Sent clipboard (%d bytes) to %s", len(data), ip)
	}
}
