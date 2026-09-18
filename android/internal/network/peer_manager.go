package network

import (
	"clipsync-android/internal"
	"clipsync-android/internal/events"
	"fmt"
	"sync"
	"time"
)

type PeerManager struct {
	mu    sync.RWMutex
	peers map[string]internal.Device
	bus   *events.EventBus
}

func NewPeerManager(bus *events.EventBus) *PeerManager {
	return &PeerManager{
		peers: make(map[string]internal.Device),
		bus:   bus,
	}
}

func (pm *PeerManager) AddOrUpdatePeer(dev internal.Device) {
	pm.mu.Lock()
	existing, found := pm.peers[dev.IP]
	isNew := !found

	pm.peers[dev.IP] = dev
	pm.mu.Unlock()

	if pm.bus != nil {
		if isNew {
			pm.bus.PublishDeviceEvent(events.DeviceEvent{
				Type: events.DeviceConnected,
				Device: events.DeviceInfo{
					Name: dev.Name,
					IP:   dev.IP,
				},
			})
			pm.bus.PublishToast(fmt.Sprintf("Connected to %s", dev.Name))
		} else if !existing.Alive && dev.Alive {
			pm.bus.PublishDeviceEvent(events.DeviceEvent{
				Type: events.DeviceConnected,
				Device: events.DeviceInfo{
					Name: dev.Name,
					IP:   dev.IP,
				},
			})
		}
	}
}

func (pm *PeerManager) TouchPeer(ip string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if dev, found := pm.peers[ip]; found {
		dev.Alive = true
		dev.LastSeen = time.Now()
		pm.peers[ip] = dev
	}
}

func (pm *PeerManager) PruneDeadPeers(timeout time.Duration) []string {
	pm.mu.Lock()
	var pruned []string
	now := time.Now()

	for ip, dev := range pm.peers {
		if now.Sub(dev.LastSeen) > timeout {
			pruned = append(pruned, ip)
			delete(pm.peers, ip)
		}
	}
	pm.mu.Unlock()

	if pm.bus != nil {
		for _, ip := range pruned {
			pm.bus.PublishDeviceEvent(events.DeviceEvent{
				Type: events.DeviceDisconnected,
				Device: events.DeviceInfo{
					IP: ip,
				},
			})
		}
	}
	return pruned
}

func (pm *PeerManager) GetActivePeers() []internal.Device {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	var list []internal.Device
	for _, dev := range pm.peers {
		if dev.Alive {
			list = append(list, dev)
		}
	}
	return list
}

func (pm *PeerManager) GetPeer(ip string) (internal.Device, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	dev, found := pm.peers[ip]
	return dev, found
}

func (pm *PeerManager) RemovePeer(ip string) {
	pm.mu.Lock()
	dev, found := pm.peers[ip]
	if found {
		delete(pm.peers, ip)
	}
	pm.mu.Unlock()

	if found && pm.bus != nil {
		pm.bus.PublishDeviceEvent(events.DeviceEvent{
			Type: events.DeviceDisconnected,
			Device: events.DeviceInfo{
				Name: dev.Name,
				IP:   dev.IP,
			},
		})
	}
}
