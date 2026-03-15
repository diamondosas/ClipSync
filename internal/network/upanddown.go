package network

import (
	// "bufio"
	// "fmt"
	"context"
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

func RecieveClipboard(ctx context.Context) ([]byte, int){
	for {
	buffer := make([]byte, 1024)
	n, _, err := Conn.ReadFromUDP(buffer)
	if err != nil{
		log.Println("Error", err)
	}
	return buffer, n
	}
}
