package network

import (
	// "bufio"
	// "fmt"
	"log"
	// sysClipboard "golang.design/x/clipboard"
	// "clipsync/internal/globals"
)
var Buffer []byte
func SendData(data []byte) {
	_, err := Conn.Write(data)
	if err != nil {
		log.Println(err)
	}
	log.Println("Sent Clipboard: ", data)
}

func RecieveClipboard() ([]byte, int){
	for {
	Buffer = make([]byte, 1024)
	n, addr, err := Conn.ReadFromUDP(Buffer)
	if err != nil{
		log.Println("Error", err)
	}
	log.Println("Recieved Clipboard From Addr: ", addr, "Content", string(Buffer[:n]))
	return Buffer, n
	}
}
