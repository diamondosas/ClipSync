package network

import (
	// "bufio"
	"log"
	"net"
	"strconv"

	"clipsync/internal/globals"
)

// type Info struct {

// 	ConnectedTo map[string]string
// 	Dialer 		bool
// }

var InConn *net.UDPConn
var OutConn *net.UDPConn
var Ln net.Listener
var Ready = make(chan struct{})
func Connect(ip string) {

	addr, err := net.ResolveUDPAddr("udp", ip + ":" + strconv.Itoa(globals.PORT))
	if err != nil {
		log.Println(err)
	}
	
	OutConn, err = net.DialUDP("udp", nil, addr)
	// defer Conn.Close()
	if err != nil {
		log.Println(err)
	}

	// SendDetails()
}

func Listen() error {
	addr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(globals.PORT))
	InConn, _ = net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Listening For Connection...")
	close(Ready)
	return nil
}


// func SendDetails(){
// 	message := []byte("$IP-ADDR:")
// 	_, err := Conn.Write(message)

// }