package network

import (
	// "bufio"
	// "fmt"
	"log"
	"net"
	// sysClipboard "golang.design/x/clipboard"
	"clipsync/internal/globals"
)
var Buffer []byte
func SendClipboard(data []byte) {
	for _, ip := range globals.IPS{
		addr, err := net.ResolveUDPAddr("udp", ip + globals.PORT)
		if err != nil {
			panic(err)
		}
			Conn.WriteToUDP([]byte(data), addr)
	}
}

func RecieveClipboard() ([]byte, int){
	Buffer = make([]byte, 1024)
	n, addr, err := Conn.ReadFromUDP(Buffer)
	if err != nil{
		log.Println("Error", err)
	}
	log.Println("Recieved Clipboard From Addr: ", addr, "Content", string(Buffer[:n]))
	return Buffer, n
}
