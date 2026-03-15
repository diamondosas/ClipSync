package main

import (

	"context"
	"fmt"
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
	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	
	eg.Go(func() error{
		return network.RegisterDevice(ctx)
	})
	
	eg.Go(func() error{
		return network.BrowseForDevices()
	})

	eg.Go(func() error{
		return network.Listen()
	})

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		changedText := clipboard.WatchClipboard(ctx) // make this return the channel
		for data := range changedText {
			if slices.Equal(data, network.Buffer) {
				continue
			} else {
				network.SendData(data)
			}
		}
		<-ctx.Done()
	}()
	go func() {

		<-network.Ready
		for {
			buffer, n := network.RecieveClipboard()
			clipboard.WriteClipboard(string(buffer[:n]))
		}
	}()

	
}

func RunWithContext(ctx context.Context, task func()) {
	
	for{
		select{
			case <-ctx.Done():
				log.Println("Shutting Down Gracefully")
				return
			default:
				go task()
		}
	}
}
