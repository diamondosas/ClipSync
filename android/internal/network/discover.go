package network

import (
	"clipsync-android/internal"
	"context"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

type Discoverer struct {
	pm       *PeerManager
	sender   *Sender
	hostname string
	port     int

	mu            sync.Mutex
	currentCancel context.CancelFunc
	lastSig       string
}

func NewDiscoverer(pm *PeerManager, sender *Sender, hostname string, port int) *Discoverer {
	if hostname == "" {
		h, err := os.Hostname()
		if err != nil || h == "" {
			hostname = "Android-Device"
		} else {
			hostname = h
		}
	}
	if port <= 0 {
		p, _ := strconv.Atoi(internal.DefaultPort)
		port = p
	}
	return &Discoverer{
		pm:       pm,
		sender:   sender,
		hostname: hostname,
		port:     port,
	}
}

func (d *Discoverer) Start(ctx context.Context) error {
	ticker := time.NewTicker(internal.DiscoveryScanInterval)
	defer ticker.Stop()

	d.restartDiscovery(ctx)

	for {
		select {
		case <-ctx.Done():
			d.mu.Lock()
			if d.currentCancel != nil {
				d.currentCancel()
			}
			d.mu.Unlock()
			return ctx.Err()
		case <-ticker.C:
			newSig := GetCurrIPS()
			d.mu.Lock()
			changed := (newSig != d.lastSig)
			d.mu.Unlock()
			if changed {
				log.Println("[Discovery] Network state changed! Re-registering...")
				d.restartDiscovery(ctx)
			}
		}
	}
}

func (d *Discoverer) restartDiscovery(ctx context.Context) {
	d.mu.Lock()
	if d.currentCancel != nil {
		d.currentCancel()
	}
	sig := GetCurrIPS()
	d.lastSig = sig

	if sig == "" {
		d.mu.Unlock()
		log.Println("[Discovery] No active network interface found, waiting...")
		return
	}

	log.Printf("[Discovery] Active network detected (%s), starting Zeroconf...", sig)

	subCtx, cancel := context.WithCancel(ctx)
	d.currentCancel = cancel
	d.mu.Unlock()

	go d.registerService(subCtx)
	go d.browseServices(subCtx)
}

func (d *Discoverer) registerService(ctx context.Context) error {
	server, err := zeroconf.Register(d.hostname, internal.ServiceType, internal.ServiceDomain, d.port, []string{""}, nil)
	if err != nil {
		log.Printf("[Discovery] Zeroconf register error: %v", err)
		return err
	}
	defer server.Shutdown()

	<-ctx.Done()
	return nil
}

func (d *Discoverer) browseServices(ctx context.Context) error {
	resolver, err := zeroconf.NewResolver()
	if err != nil {
		log.Printf("[Discovery] Zeroconf resolver error: %v", err)
		return err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go d.handleEntries(entries)

	err = resolver.Browse(ctx, internal.ServiceType, internal.ServiceDomain, entries)
	if err != nil {
		log.Printf("[Discovery] Zeroconf browse error: %v", err)
		return err
	}

	<-ctx.Done()
	return nil
}

func (d *Discoverer) handleEntries(entries <-chan *zeroconf.ServiceEntry) {
	for entry := range entries {
		if entry.Instance != d.hostname && len(entry.AddrIPv4) > 0 {
			newIP := entry.AddrIPv4[0].String()
			name := entry.Instance
			if name == "" {
				name = entry.HostName
			}

			if d.pm != nil {
				d.pm.AddOrUpdatePeer(internal.Device{
					Name:     name,
					IP:       newIP,
					Alive:    true,
					LastSeen: time.Now(),
				})
			}

			if d.sender != nil {
				_ = d.sender.SendHandshake(newIP, d.hostname)
			}

			log.Printf("[Discovery] Discovered ClipSync peer: %s (%s)", name, newIP)
		}
	}
}

// ConnectManual allows connecting directly to a device IP when mDNS is unavailable.
func (d *Discoverer) ConnectManual(ip string) error {
	if d.pm != nil {
		d.pm.AddOrUpdatePeer(internal.Device{
			Name:     ip,
			IP:       ip,
			Alive:    true,
			LastSeen: time.Now(),
		})
	}
	if d.sender != nil {
		return d.sender.SendHandshake(ip, d.hostname)
	}
	return nil
}

// GetCurrIPS returns a sorted, comma-separated list of active non-loopback IPv4 addresses.
func GetCurrIPS() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	var ips []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ips = append(ips, ipNet.IP.String())
				}
			}
		}
	}
	sort.Strings(ips)
	return strings.Join(ips, ",")
}
