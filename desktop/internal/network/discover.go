package network

import (
	"clipsync/internal"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	"strings"
	"sort"
	
	"clipsync/internal/view"

	"github.com/grandcat/zeroconf"
)

var currentCancel context.CancelFunc
var lastsig string

func StartAutoDiscovery(ctx context.Context) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	restartDiscovery(ctx)

	for {
		select {
		case <-ctx.Done():
			if currentCancel != nil {
				currentCancel()
			}
			return ctx.Err()
		case <-ticker.C:
			newSig := GetCurrIPS()
			if newSig != lastsig {
				log.Println("[Discovery] Network state changed! Re-registering...")
				restartDiscovery(ctx)
			}
		}
	}
}

func restartDiscovery(ctx context.Context) {
	if currentCancel != nil {
		currentCancel()
	}
	sig := GetCurrIPS()

	lastsig = sig

	if sig == "" {
		log.Println("[Discovery] No active network found, waiting...")
		return
	}

	log.Printf("[Discovery] Network change detected (%s), starting Zeroconf...", sig)

	var subCtx context.Context
	subCtx, cancel := context.WithCancel(ctx)
	currentCancel = cancel

	go RegisterDevice(subCtx)
	go BrowseForDevices(subCtx)
}

func RegisterDevice(ctx context.Context) error {
	internal.Hostname, _ = os.Hostname()
	name := internal.Hostname

	intPORT, _ := strconv.Atoi(internal.PORT)

	server, err := zeroconf.Register(name, "_clipsync._tcp", "local.", int(intPORT), []string{""}, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Broadcasting Presence...")
	defer server.Shutdown()

	<-ctx.Done()
	return nil
}

// Discover all services on the network (e.g. _workstation._tcp)

func BrowseForDevices(ctx context.Context) error {
	resolver, err := zeroconf.NewResolver()

	if err != nil {
		log.Println(err)
		return err
	}
	var entries = make(chan *zeroconf.ServiceEntry)

	go entry(entries)

	err = resolver.Browse(ctx, "_clipsync._tcp", "local.", entries)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Starting to Discover Services...")

	<-ctx.Done()

	return nil
}

func entry(results <-chan *zeroconf.ServiceEntry) {
	for entry := range results {
		if entry.Instance != internal.Hostname && len(entry.AddrIPv4) > 0 {
			newIP := entry.AddrIPv4[0].String()
			newDevice := internal.Device{Name: entry.HostName, Ip: newIP, Alive: true, LastSeen: time.Now()}

			
			view.AddNewDevice(newDevice)
			Connect(newIP)

			log.Println("Found Device: Name: ", entry.Instance, " IP: ", entry.AddrIPv4)
			fmt.Println("Connected Device:", entry.Instance)
		}
	}
}




// GetNetworkSignature returns a string signature of active IPv4 addresses
func GetCurrIPS() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	var ips []string
	for _, iface := range ifaces {
		// Only check interfaces that are UP and not Loopback
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
