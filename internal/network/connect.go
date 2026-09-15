package network

import (
	"clipsync/internal"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/xtaci/kcp-go/v5"
)

var(
	Sess *kcp.UDPSession
	SessMu sync.Mutex
)

var Ready = make(chan struct{})
var	BlockCrypt kcp.BlockCrypt

func Connect(ip string) {
	if Sess == nil {
		log.Println("Cannot connect, Conn is not initialized. Waiting for Ready channel...")
		<-Ready
	}

	block := prepareCrypt()

	Sess, err := kcp.DialWithOptions(ip + ":" + internal.PORT, block, 10, 3)
	if err != nil{
		log.Println("Could not Connect to ", ip)
	}
	defer Sess.Close()

	//Ultra low latency mode
	Sess.SetNoDelay(1, 10, 2, 1)

	msg := []byte{MsgTypeHandshake}

	_, err = Sess.Write(msg)
	if err != nil {
		log.Println("Connect Write error:", err)
	}
}

func Listen(ctx context.Context) error {

	block := prepareCrypt()
	listener, err := kcp.ListenWithOptions("0.0.0.0:" + internal.PORT, block, 10, 3)
	if err != nil{
		log.Println("Could not Bind on Port ", internal.PORT)
		log.Fatal(err)
	}
	close(Ready)
	defer listener.Close()
	go func(){
		for{
			sess, err := listener.AcceptKCP()
			if err != nil{
				log.Println("Could not Accept connection")
			}
			msg := []byte{MsgTypeHandshake}
			_, err = sess.Write(msg)
			if err != nil{
				log.Println("Could not Send Handshake", err)
			}
		}
	}()
	fmt.Println("Udp server listeining on port: ", internal.PORT)


	<-ctx.Done()
	if listener != nil{
		listener.Close()
	}
	
	return nil
}

func prepareCrypt() kcp.BlockCrypt {
	block, err := kcp.NewAESBlockCrypt(internal.SecretKey)
	if err != nil{
		log.Println("Could not create block")
	}
	if BlockCrypt == nil{
		BlockCrypt = block
	}
	return block
}