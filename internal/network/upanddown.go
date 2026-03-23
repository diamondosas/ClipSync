package network

import (
	// "bufio"
	// "fmt"
	"log"
	"net"
	"strconv"
	"slices"

	// sysClipboard "golang.design/x/clipboard"
	"clipsync/internal/globals"
)
var Buffer []byte
func SendClipboard(data []byte) {
	if Conn == nil {
		log.Println("SendClipboard: Conn is nil, skipping.")
		return
	}
	
	globals.IPSMu.Lock()
	ips := make([]string, len(globals.IPS))
	copy(ips, globals.IPS)
	globals.IPSMu.Unlock()

	for _, ip := range ips {
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
	if slices.Equal(Buffer, []byte("---ClipSync---")){
		globals.IPSMu.Lock()
		globals.IPS = append(globals.IPS, string(addr.IP))
		globals.IPSMu.Unlock()
	}else{
		log.Println("Recieved Clipboard From Addr: ", addr, "Content", string(Buffer[:n]))
		return Buffer, n
	}

	return nil, 0
}
