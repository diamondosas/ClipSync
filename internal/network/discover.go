package network

import (
	"context"
	"fmt"
	"log"
	"os"

	"clipsync/internal/globals"

	"github.com/grandcat/zeroconf"
)

var Entries = make(chan *zeroconf.ServiceEntry)

func RegisterDevice(ctx context.Context, name string) error {
	if name == "" {
		globals.Username, _ = os.Hostname()
		name = globals.Username
	}

	server, err := zeroconf.Register(name, "_clipsync._tcp", "local.", globals.PORT, []string{""}, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Broadcasting Presence...")
	defer server.Shutdown()
	select {
	case <-ctx.Done():
		return nil
	}
}

// Discover all services on the network (e.g. _workstation._tcp)

func BrowseForDevices(ctx context.Context) error {
	reslover, err := zeroconf.NewResolver(nil)

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
	select {
	case <-ctx.Done():
		return nil
	}
}

func entry(results <-chan *zeroconf.ServiceEntry) {
	for entry := range results {
		if entry.Instance != globals.Username {
			newIP := string(entry.AddrIPv4[0].String())
			globals.IPS = append(globals.IPS, newIP)
			// Connect(newIP)
			log.Println("Found Device: Name: ", entry.Instance, " IP: ", entry.AddrIPv4)

			fmt.Println("Connected Device:", entry.Instance)
		}
	}
}
