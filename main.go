package main

import (

	"context"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"
	
	"clipsync/internal/clipboard"
	"clipsync/internal/network"
	"golang.org/x/sync/errgroup"
	// "clipsync/internal/ping"
)

var Version = "dev"

func main() {

	clipboard.Init()
	
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	
	eg.Go(func() error{
		return network.RegisterDevice(ctx, "")
	})
	
	eg.Go(func() error{
		return network.BrowseForDevices(ctx)
	})

	eg.Go(func() error{
		return network.Listen()
	})

	eg.Go(func() error{
		for {
			data := clipboard.WatchClipboard(ctx)
			if !slices.Equal(data, network.Buffer) {
				network.SendData(data)
			}
		}
	})
	// eg.Go(func() error {
	// 	<-network.Ready
	// 	for {
	// 		buffer, n := network.RecieveClipboard()
	// 		clipboard.WriteClipboard(string(buffer[:n]))
	// 	}
	// })
	err := eg.Wait()
	if err != nil{
		log.Fatal("Shutdown Error", err)
	}
	
}

