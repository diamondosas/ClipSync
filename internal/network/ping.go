package network

import (
	"encoding/binary"
	"log"
	"net"
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

	addr, err := net.ResolveUDPAddr("udp", ip+":"+PORT)

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

}
