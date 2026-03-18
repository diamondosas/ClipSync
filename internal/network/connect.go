package network

import (
	// "bufio"
	"context"
	"log"
	"net"
	"strconv"

	"clipsync/internal/globals"
)

// type Info struct {

// 	ConnectedTo map[string]string
// 	Dialer 		bool
// }

var Conn *net.UDPConn
var Ready = make(chan struct{})
func Connect(ip string, addr  ...*net.UDPConn) {
	addr, err := net.ResolveUDPAddr("udp", ip + ":" + strconv.Itoa(globals.PORT))
	if err != nil {
		log.Println(err)
	}
	SendClipboard([]byte("V"))
}

func Listen(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(globals.PORT))
	Conn, _ = net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err)
	}
	
	log.Println("Listening For Connection...")
	close(Ready)
	
	select{
	case <-ctx.Done():
		return nil
	}
}


// func SendDetails(){
// 	message := []byte("$IP-ADDR:")
// 	_, err := Conn.Write(message)

// }