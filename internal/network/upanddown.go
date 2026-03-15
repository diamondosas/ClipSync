package network

import (
	// "bufio"
	// "fmt"
	"log"
	// sysClipboard "golang.design/x/clipboard"
	// "clipsync/internal/globals"
)

func SendData(data []byte) {
	_, err := Conn.Write(data)
	if err != nil {
		log.Println(err)
	}
	log.Println("Sent Clipboard: ", data)
}

func RecieveClipboard() ([]byte, int){
	for {
	buffer := make([]byte, 1024)
	n, addr, err := Conn.ReadFromUDP(buffer)
	if err != nil{
		log.Println("Error", err)
	}
	log.Println("Recieved Clipboard From Addr: ", addr, "Content", string(buffer[:n]))
	return buffer, n
	}
}
