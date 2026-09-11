package network

import (
	"clipsync/internal"
	"encoding/binary"
	"log"
	"net"
	"time"
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
	if Conn == nil {
		return
	}

	addr, err := net.ResolveUDPAddr("udp", ip + ":" + internal.PORT)

	if err != nil {
		log.Println("SendPing Resolve Error:", err)
		return
	}

	msg := []byte("~")
	payload := make([]byte, 4+len(msg))
	binary.BigEndian.PutUint32(payload[:4], uint32(len(msg)))
	copy(payload[4:], msg)

	_, err = Conn.WriteToUDP(payload, addr)

	if err != nil {
		log.Println("SendPing Write Error:", err)
	}
}

func CheckForPing() {
	internal.ConnDevicesMu.Lock()
	time.Sleep(time.Second * 1)
	for i := range 	internal.ConnDevices{
		if string(Addr.IP.String()) == string(internal.ConnDevices[i].Ip){
			if (time.Since(internal.ConnDevices[i].LastSeen) <= (time.Second * 11)){
				internal.ConnDevices[i].Alive = false
			}
		}
	}
	internal.ConnDevicesMu.Unlock()
}
