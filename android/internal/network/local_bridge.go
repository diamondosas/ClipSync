package network

import (
	"clipsync-android/internal/events"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

const DefaultLocalBridgePort = 9998

// LocalBridge listens on localhost UDP for copied text emitted by the Android AccessibilityService.
type LocalBridge struct {
	bus      *events.EventBus
	pm       *PeerManager
	sender   *Sender
	port     int
	lastClip string
	mu       sync.Mutex
}

func NewLocalBridge(bus *events.EventBus, pm *PeerManager, sender *Sender, port int) *LocalBridge {
	if port <= 0 {
		port = DefaultLocalBridgePort
	}
	return &LocalBridge{
		bus:    bus,
		pm:     pm,
		sender: sender,
		port:   port,
	}
}

func (b *LocalBridge) SetSender(sender *Sender) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sender = sender
}

func (b *LocalBridge) Start(ctx context.Context) (string, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", b.port))
	if err != nil {
		return "", err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return "", fmt.Errorf("local bridge bind failed: %w", err)
	}

	boundPort := strconv.Itoa(conn.LocalAddr().(*net.UDPAddr).Port)
	log.Printf("[LocalBridge] Android Accessibility IPC listening on 127.0.0.1:%s", boundPort)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	go func() {
		buf := make([]byte, 65535)
		for {
			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}

			if n <= 0 {
				continue
			}

			text := string(buf[:n])

			b.mu.Lock()
			if text == b.lastClip {
				b.mu.Unlock()
				continue
			}
			b.lastClip = text
			currentSender := b.sender
			b.mu.Unlock()

			// 1. Dispatch locally to UI event bus
			if b.bus != nil {
				b.bus.PublishClipReceived(events.ClipReceivedEvent{
					Content: text,
					FromIP:  "127.0.0.1 (Phone)",
				})
				b.bus.PublishToast("Phone copy synced to devices")
			}

			// 2. Broadcast to all connected desktop/remote peers
			if currentSender != nil {
				currentSender.BroadcastClipboard([]byte(text))
			}
		}
	}()

	return boundPort, nil
}
