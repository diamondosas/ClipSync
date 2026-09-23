package network

import (
	"clipsync/internal"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var (
	serverConn   *net.UDPConn
	serverConnMu sync.RWMutex
)

// SendPacket sends a UDP packet to the given IP on the default port.
func SendPacket(ip string, data []byte) error {
	targetAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(ip, internal.PORT))
	if err != nil {
		return err
	}

	serverConnMu.RLock()
	conn := serverConn
	serverConnMu.RUnlock()

	if conn != nil {
		_, err = conn.WriteToUDP(data, targetAddr)
		return err
	}

	// Fallback to a temporary UDP socket if listener is not ready
	c, err := net.DialUDP("udp", nil, targetAddr)
	if err != nil {
		return err
	}
	defer c.Close()
	_, err = c.Write(data)
	return err
}

func Connect(ip string) {
	if internal.Hostname == "" {
		internal.Hostname, _ = os.Hostname()
	}
	msg := append([]byte{MsgTypeHandshake}, []byte(internal.Hostname)...)

	err := SendPacket(ip, msg)
	if err != nil {
		log.Println("Connect Write error:", err)
		return
	}
	log.Printf("[Network] Sent Handshake to %s", ip)
}

func Listen(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+internal.PORT)
	if err != nil {
		log.Println("Could not resolve UDP address for Port", internal.PORT)
		return err
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("Could not Bind on Port", internal.PORT)
		return err
	}
	defer listener.Close()

	serverConnMu.Lock()
	serverConn = listener
	serverConnMu.Unlock()

	fmt.Println("UDP server listening on port:", internal.PORT)

	go func() {
		buf := make([]byte, 65535)
		for {
			n, remoteAddr, err := listener.ReadFromUDP(buf)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					log.Println("Could not read from UDP:", err)
					continue
				}
			}
			if n <= 0 {
				continue
			}

			packet := make([]byte, n)
			copy(packet, buf[:n])

			go HandleIncomingPacket(packet, remoteAddr)
		}
	}()

	<-ctx.Done()

	serverConnMu.Lock()
	if serverConn != nil {
		serverConn.Close()
		serverConn = nil
	}
	serverConnMu.Unlock()

	return nil
}