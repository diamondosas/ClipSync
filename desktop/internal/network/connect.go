package network

import (
	"clipsync/internal"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/xtaci/kcp-go/v5"
)

var(
	peers = make(map[string]*kcp.UDPSession)
	peersMu sync.Mutex

	BlockCrypt kcp.BlockCrypt
	// Ready = make(chan struct{})
)

func Connect(ip string) {
	sess, err := createPeer(ip)
	if err != nil{
		log.Println("Could not connect to ip:", ip)
		return
	}

	if internal.Hostname == "" {
		internal.Hostname, _ = os.Hostname()
	}
	msg := append([]byte{MsgTypeHandshake}, []byte(internal.Hostname)...)

	_, err = sess.Write(msg)
	log.Println("Sent HandShake")
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
	defer listener.Close()

	fmt.Println("Udp server listeining on port: ", internal.PORT)


	go func(){
		for{
			sess, err := listener.AcceptKCP()
			if err != nil{
				select{
				case <-ctx.Done():
					return 
				
				default:
					log.Println("Could not Accept connection")
					continue
				}
			}
			sess.SetNoDelay(1, 10, 2, 1)

			go HandleIncomingData(sess)		
		}
	}()


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

func createPeer(ip string) (*kcp.UDPSession, error){
	peersMu.Lock()
	sess, exists := peers[ip]
	peersMu.Unlock()

	if exists && sess != nil{
		return sess, nil
	}

	block := prepareCrypt()

	sess, err := kcp.DialWithOptions(ip + ":" + internal.PORT, block, 10, 3)
	if err != nil{
		log.Println("Could not Connect to ", ip)
		return nil, fmt.Errorf("Could not conect to %s", ip) 
	}
	peersMu.Lock()
	peers[ip] = sess
	peersMu.Unlock()
	
	sess.SetNoDelay(1, 10, 2, 1)

	return sess, nil
}

func removePeer(ip string){
	peersMu.Lock()
	defer peersMu.Unlock()
	if sess, exists := peers[ip]; exists {
		if sess != nil{
			sess.Close()
		}
		delete(peers, ip)
	}
}