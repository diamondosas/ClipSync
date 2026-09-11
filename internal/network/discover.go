package network

import (
	"clipsync/internal"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/grandcat/zeroconf"
	"clipsync/internal/view"
)

var Entries = make(chan *zeroconf.ServiceEntry)

func RegisterDevice(ctx context.Context) error {
	internal.Username, _ = os.Hostname()
	name := internal.Username

	ifaces := getAllInterfaces()

	intPORT, _ := strconv.ParseInt(internal.PORT, 10, 8)

	server, err := zeroconf.Register(name, "_clipsync._tcp", "local.", int(intPORT), []string{""}, ifaces)

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
	ifaces := getAllInterfaces()
	reslover, err := zeroconf.NewResolver(zeroconf.SelectIfaces(ifaces))

	if err != nil {
		log.Println(err)
		return err
	}

	go entry(Entries)

	err = reslover.Browse(ctx, "_clipsync._tcp", "local.", Entries)

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
		if entry.Instance != internal.Username && len(entry.AddrIPv4) > 0 {
			newIP := entry.AddrIPv4[0].String()
			newDevice := internal.Device{Name: entry.HostName, Ip: newIP, Alive: true}

			internal.IPSMu.Lock()
			exists := false
			for _, ip := range internal.IPS {
				if ip == newIP {
					exists = true
					break
				}
			}
			if !exists {
				internal.IPS = append(internal.IPS, newIP)
			}
			internal.IPSMu.Unlock()
			go view.AddNewDevice(newDevice)
			go Connect(newIP)
			log.Println("Found Device: Name: ", entry.Instance, " IP: ", entry.AddrIPv4)

			fmt.Println("Connected Device:", entry.Instance)
		}
	}
}

func getAllInterfaces() []net.Interface {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var result []net.Interface
	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// Skip interfaces with no addresses
		addrs, err := iface.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}
		result = append(result, iface)
	}
	return result
}
