package network

import (
	"context"
	"fmt"
	"log"
	"os"

	"clipsync/internal/globals"

	"github.com/grandcat/zeroconf"
	"golang.org/x/text/cases"
)

var Entries = make(chan *zeroconf.ServiceEntry)


func RegisterDevice(ctx context.Context) error {
	globals.Username, _ = os.Hostname()
	
	server, err := zeroconf.Register(globals.Username, "_clipsync._tcp", "local.", globals.PORT, []string{""}, nil)
	
	if err != nil {
		log.Println(err)
		return err
	}
	
	log.Println("Broadcasting Presence...")
	defer server.Shutdown()
	select{
	case <-ctx.Done():
		return nil
	}
}

// Discover all services on the network (e.g. _workstation._tcp)

func BrowseForDevices() error{
	reslover, err := zeroconf.NewResolver(nil)
	
	if err != nil {
		log.Println(err)
		return err
	}

	go entry(Entries)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	err = reslover.Browse(ctx, "_clipsync._tcp", "local.", Entries)
	
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Starting to Discover Services")

}

func entry(results <-chan *zeroconf.ServiceEntry) {
	for entry:= range results{
		if entry.Instance != globals.Username {
			// ip := string(entry.AddrIPv4[0].String())
			// Connect(ip)
			log.Println("Found Device: Name: ", entry.Instance," IP: ", entry.AddrIPv4)
			globals.IP = append(globals.IP, string(entry.AddrIPv4[0].String()))
			fmt.Println("Connected Device:", entry.Instance)
		}
	}
}

