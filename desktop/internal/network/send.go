package network

import (
	"clipsync/internal"
	"log"

)

//Sends Clipboard To all Connected Peers
func SendClipboard(data []byte) {
	var ips []string

	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices{
		ips = append(ips, internal.ConnDevices[i].Ip)
	} 
	internal.ConnDevicesMu.Unlock()
		
	payload := append([]byte{MsgTypeClipboard}, data...)

	for _, ip := range ips {
		sess , err := createPeer(ip)
		if err != nil || sess == nil{
			continue
		}
		_, err = sess.Write(payload)
		if err != nil {
			log.Println("SendClipboard Write Error:", err)
			removePeer(ip)
			continue
		}

	}
}
