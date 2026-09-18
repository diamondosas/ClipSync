package network

import (
	"clipsync-android/internal"
	"clipsync-android/internal/crypto"
	"clipsync-android/internal/events"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

type Listener struct {
	pm        *PeerManager
	bus       *events.EventBus
	secretKey []byte
	port      int
}

func NewListener(pm *PeerManager, bus *events.EventBus, secretKey []byte, port int) *Listener {
	return &Listener{
		pm:        pm,
		bus:       bus,
		secretKey: secretKey,
		port:      port,
	}
}

func (l *Listener) Start(ctx context.Context) (string, error) {
	block, err := crypto.NewBlockCrypt(l.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create crypt: %w", err)
	}

	listenAddr := fmt.Sprintf("0.0.0.0:%d", l.port)
	listener, err := kcp.ListenWithOptions(listenAddr, block, 10, 3)
	if err != nil {
		return "", fmt.Errorf("failed to bind on %s: %w", listenAddr, err)
	}

	boundPort := strconv.Itoa(listener.Addr().(*net.UDPAddr).Port)
	log.Printf("[Listener] UDP/KCP listening on %s (port %s)", listenAddr, boundPort)

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	go func() {
		for {
			sess, err := listener.AcceptKCP()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}
			sess.SetNoDelay(1, 10, 2, 1)
			go l.handleIncomingSession(ctx, sess)
		}
	}()

	return boundPort, nil
}

func (l *Listener) handleIncomingSession(ctx context.Context, sess *kcp.UDPSession) {
	defer sess.Close()

	remoteAddr := sess.RemoteAddr().(*net.UDPAddr)
	ip := remoteAddr.IP.String()
	tmpBuf := make([]byte, 65535)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := sess.Read(tmpBuf)
		if err != nil {
			return
		}

		msgType, payload, err := DecodePacket(tmpBuf[:n])
		if err != nil {
			log.Printf("[Listener] Malformed packet from %s: %v", ip, err)
			continue
		}

		switch msgType {
		case internal.MsgTypeHandshake:
			hostname := string(payload)
			if hostname == "" {
				hostname = ip
			}
			if l.pm != nil {
				l.pm.AddOrUpdatePeer(internal.Device{
					Name:     hostname,
					IP:       ip,
					Alive:    true,
					LastSeen: time.Now(),
				})
			}

		case internal.MsgTypePing:
			if l.pm != nil {
				l.pm.TouchPeer(ip)
			}

		case internal.MsgTypeClipboard:
			content := string(payload)
			if len(content) > 0 {
				if l.bus != nil {
					l.bus.PublishClipReceived(events.ClipReceivedEvent{
						Content: content,
						FromIP:  ip,
					})
					l.bus.PublishToast("Synced new clip")
				}
			}
		}
	}
}
