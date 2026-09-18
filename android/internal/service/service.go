package service

import (
	"clipsync-android/internal"
	"clipsync-android/internal/events"
	"clipsync-android/internal/network"
	"clipsync-android/internal/storage"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

type Config struct {
	Port      int
	Hostname  string
	SecretKey []byte
}

type ClipSyncService struct {
	bus        *events.EventBus
	config     Config
	pm         *network.PeerManager
	sender     *network.Sender
	listener   *network.Listener
	discoverer  *network.Discoverer
	localBridge *network.LocalBridge
	storage     *storage.Storage
}

func NewClipSyncService(bus *events.EventBus, config Config) *ClipSyncService {
	if bus == nil {
		bus = events.NewEventBus()
	}
	if config.Hostname == "" {
		h, err := os.Hostname()
		if err != nil || h == "" {
			config.Hostname = "Android-ClipSync"
		} else {
			config.Hostname = h
		}
	}
	if config.Port <= 0 {
		p, _ := strconv.Atoi(internal.DefaultPort)
		config.Port = p
	}
	if len(config.SecretKey) == 0 {
		config.SecretKey = internal.SecretKey
	}

	pm := network.NewPeerManager(bus)
	sender := network.NewSender(pm, config.SecretKey, config.Port)
	listener := network.NewListener(pm, bus, config.SecretKey, config.Port)
	discoverer := network.NewDiscoverer(pm, sender, config.Hostname, config.Port)
	localBridge := network.NewLocalBridge(bus, pm, sender, 0)
	store := storage.NewStorage()

	return &ClipSyncService{
		bus:         bus,
		config:      config,
		pm:          pm,
		sender:      sender,
		listener:    listener,
		discoverer:  discoverer,
		localBridge: localBridge,
		storage:     store,
	}
}

func (s *ClipSyncService) Start(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	// 1. Start KCP Listener
	boundPortStr, err := s.listener.Start(egCtx)
	if err != nil {
		return err
	}
	if boundPort, err := strconv.Atoi(boundPortStr); err == nil && boundPort > 0 {
		// Update sender & discoverer port if dynamic
		s.discoverer = network.NewDiscoverer(s.pm, s.sender, s.config.Hostname, boundPort)
	}

	// 2. Start LocalBridge for Android Accessibility Service IPC
	if _, err := s.localBridge.Start(egCtx); err != nil {
		log.Printf("[ClipSyncService] LocalBridge warning: %v", err)
	}

	// 3. Start Auto Discovery
	eg.Go(func() error {
		return s.discoverer.Start(egCtx)
	})

	// 3. Heartbeat Ping Loop (every 2s)
	eg.Go(func() error {
		ticker := time.NewTicker(internal.PingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-ticker.C:
				s.sender.PingAllPeers()
			}
		}
	})

	// 4. Dead Peer Prune Loop (every 1s)
	eg.Go(func() error {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-ticker.C:
				s.pm.PruneDeadPeers(internal.PeerTimeout)
			}
		}
	})

	log.Println("[ClipSyncService] Background service started successfully")
	err = eg.Wait()
	s.sender.CloseAll()
	return err
}

func (s *ClipSyncService) BroadcastClip(text string) {
	if text == "" {
		return
	}
	s.sender.BroadcastClipboard([]byte(text))
}

func (s *ClipSyncService) ConnectIP(ip string) error {
	return s.discoverer.ConnectManual(ip)
}

func (s *ClipSyncService) GetConnectedDevices() []internal.Device {
	return s.pm.GetActivePeers()
}

func (s *ClipSyncService) PeerManager() *network.PeerManager {
	return s.pm
}

func (s *ClipSyncService) EventBus() *events.EventBus {
	return s.bus
}

func (s *ClipSyncService) Storage() *storage.Storage {
	return s.storage
}
