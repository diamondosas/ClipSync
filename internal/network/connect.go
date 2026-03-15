package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"

	"clipsync/internal/globals"
	"clipsync/internal/clipboard"
)

// type Info struct {

// 	ConnectedTo map[string]string
// 	Dialer 		bool
// }

var Conn net.Conn
var Ln net.Listener

func Connect(ip string) {

	addr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		log.Println(err)
	}

	//Send and recive confirm form server
	Conn, err := net.DialUDP("udp", nil, addr)
	defer Conn.Close()
	if err != nil {
		log.Println(err)
	}

	SendDetails()
}

func Listen() {
	defer globals.WG.Done()
	addr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(globals.PORT))
	Conn, _ := net.ListenUDP("udp", addr)
	defer Conn.Close()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Listening...")
	for {
		buffer := make([]byte, 1024)
		n, _, err := Conn.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("Error", err)
		}

		WriteClipboard(string(buffer[:n]))
		}


}


func SendDetails(){
	message := []byte("")
	Conn.Write(message)
}