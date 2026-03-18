package network

import (
	// "bufio"
	// "fmt"
	"log"
	"net"
	"strconv"

	// sysClipboard "golang.design/x/clipboard"
	"clipsync/internal/globals"
)
var Buffer []byte
func SendClipboard(data []byte) {
	if Conn == nil {
		log.Println("SendClipboard: Conn is nil, skipping.")
		return
	}
	for _, ip := range globals.IPS{
		addr, err := net.ResolveUDPAddr("udp", ip + ":" + strconv.Itoa(globals.PORT))
		if err != nil {
			log.Println("SendClipboard Resolve Error:", err)
			continue
		}
		_, err = Conn.WriteToUDP(data, addr)
		if err != nil {
			log.Println("SendClipboard Write Error:", err)
		}
	}
}

func RecieveClipboard() ([]byte, int){
	if Conn == nil {
		log.Println("RecieveClipboard: Conn is nil. Waiting for Ready...")
		<-Ready
	}
	Buffer = make([]byte, 1024)
	n, addr, err := Conn.ReadFromUDP(Buffer)
	if err != nil{
		log.Println("Error", err)
	}
	log.Println("Recieved Clipboard From Addr: ", addr, "Content", string(Buffer[:n]))
	return Buffer, n
}
