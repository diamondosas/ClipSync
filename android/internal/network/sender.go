package network

import (
	"clipsync-android/internal"
	"clipsync-android/internal/crypto"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/xtaci/kcp-go/v5"
)

type Sender struct {
	pm        *PeerManager
	secretKey []byte
	port      int

	sessionsMu sync.Mutex
	sessions   map[string]*kcp.UDPSession
}

func NewSender(pm *PeerManager, secretKey []byte, port int) *Sender {
	if port <= 0 {
		p, _ := strconv.Atoi(internal.DefaultPort)
		port = p
	}
	return &Sender{
		pm:        pm,
		secretKey: secretKey,
		port:      port,
		sessions:  make(map[string]*kcp.UDPSession),
	}
}

func (s *Sender) getOrCreateSession(ip string) (*kcp.UDPSession, error) {
	s.sessionsMu.Lock()
	sess, exists := s.sessions[ip]
	s.sessionsMu.Unlock()

	if exists && sess != nil {
		return sess, nil
	}

	block, err := crypto.NewBlockCrypt(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypt: %w", err)
	}

	targetAddr := net.JoinHostPort(ip, strconv.Itoa(s.port))
	newSess, err := kcp.DialWithOptions(targetAddr, block, 10, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", targetAddr, err)
	}

	newSess.SetNoDelay(1, 10, 2, 1)

	s.sessionsMu.Lock()
	s.sessions[ip] = newSess
	s.sessionsMu.Unlock()

	return newSess, nil
}

func (s *Sender) removeSession(ip string) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	if sess, exists := s.sessions[ip]; exists {
		if sess != nil {
			sess.Close()
		}
		delete(s.sessions, ip)
	}
}

func (s *Sender) SendHandshake(ip string, hostname string) error {
	sess, err := s.getOrCreateSession(ip)
	if err != nil {
		return err
	}

	msg := EncodeHandshake(hostname)
	_, err = sess.Write(msg)
	if err != nil {
		s.removeSession(ip)
		return err
	}
	return nil
}

func (s *Sender) SendPing(ip string) error {
	sess, err := s.getOrCreateSession(ip)
	if err != nil {
		return err
	}

	msg := EncodePing()
	_, err = sess.Write(msg)
	if err != nil {
		s.removeSession(ip)
		return err
	}
	return nil
}

func (s *Sender) PingAllPeers() {
	if s.pm == nil {
		return
	}
	peers := s.pm.GetActivePeers()
	for _, peer := range peers {
		_ = s.SendPing(peer.IP)
	}
}

func (s *Sender) BroadcastClipboard(data []byte) {
	if s.pm == nil || len(data) == 0 {
		return
	}

	peers := s.pm.GetActivePeers()
	msg := EncodeClipboard(data)

	for _, peer := range peers {
		sess, err := s.getOrCreateSession(peer.IP)
		if err != nil {
			continue
		}
		_, err = sess.Write(msg)
		if err != nil {
			s.removeSession(peer.IP)
		}
	}
}

func (s *Sender) CloseAll() {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	for ip, sess := range s.sessions {
		if sess != nil {
			sess.Close()
		}
		delete(s.sessions, ip)
	}
}
