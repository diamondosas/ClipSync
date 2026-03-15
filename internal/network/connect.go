package network

import (
	// "bufio"
	"fmt"
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
var Ln net.Listener

func Connect(ip string) {

	addr, err := net.ResolveUDPAddr("udp", ip + ":9000")
	if err != nil {
		log.Println(err)
	}

	//Send and recive confirm form server
	Conn, err = net.DialUDP("udp", nil, addr)
	// defer Conn.Close()
	if err != nil {
		log.Println(err)
	}

	// SendDetails()
}

func Listen() {
	defer globals.WG.Done()
	addr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(globals.PORT))
	Conn, _ = net.ListenUDP("udp", addr)
	defer Conn.Close()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Listening...")

}


// func SendDetails(){
// 	message := []byte("$IP-ADDR:")
// 	_, err := Conn.Write(message)

// }