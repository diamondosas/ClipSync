package network

import (
	"clipsync/internal/globals"
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

func PingIPS(ips []string) []string {
	if len(ips) == 0 {
		return nil
	}


}

func SendPing(ip string) {
	if Conn == nil {
		return
	}
	addr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(globals.PORT))
	if err != nil {
		log.Println("SendPing Resolve Error:", err)
		return
	}
	msg := []byte("---Ping---")
	payload := make([]byte, 4+len(msg))
	binary.BigEndian.PutUint32(payload[:4], uint32(len(msg)))
	copy(payload[4:], msg)
	_, err = Conn.WriteToUDP(payload, addr)
	if err != nil {
		log.Println("SendPing Write Error:", err)
	}
}

func SendPong(addr *net.UDPAddr) {
	if Conn == nil {
		return
	}
	msg := []byte("---Pong---")
	payload := make([]byte, 4+len(msg))
	binary.BigEndian.PutUint32(payload[:4], uint32(len(msg)))
	copy(payload[4:], msg)
	_, err := Conn.WriteToUDP(payload, addr)
	if err != nil {
		log.Println("SendPong Write Error:", err)
	}
}
