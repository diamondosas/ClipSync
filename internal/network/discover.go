package network

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"clipsync/internal/globals"
	"github.com/grandcat/zeroconf"
)

var Instance string
var Entries = make(chan *zeroconf.ServiceEntry)

// Add that when it display all the interfaces
// Make it to work on a perfect LAN Peer to Peer Setup

func RegisterDevice() {
	globals.Username, _ = os.Hostname()
	
	_, err := zeroconf.Register(globals.Username, "_clipsync._tcp", "local.", globals.PORT, []string{""}, nil)
	
	if err != nil {
		log.Println(err)
	}
	
	log.Println("Deivce Registered")
	// defer server.Shutdown()
}

// Discover all services on the network (e.g. _workstation._tcp)

func BrowseForDevices() {
	reslover, err := zeroconf.NewResolver(nil)
	
	if err != nil {
		log.Println(err)
	}

	go entry(Entries)
	time, cancel := context.WithTimeout(context.Background(), time.Hour*100)

	defer cancel()

	err = reslover.Browse(time, "_clipsync._tcp", "local.", Entries)
	
	if err != nil {
		log.Println(err)
	}

	log.Println("Starting to Discover Services")

}

func entry(results <-chan *zeroconf.ServiceEntry) {
	for{
		entry := <-results
		if entry == nil{
			return
		}
		if entry.Instance == globals.Username {
			continue
		} else {
			ip := string(entry.AddrIPv4[0].String())
			Connect(ip)

			log.Println("Found Device: ", entry.Instance, entry.AddrIPv4)
			globals.IP = append(globals.IP, string(entry.AddrIPv4[0].String()))
			fmt.Println("Connected Device:", entry.Instance)

		}
	}
}

