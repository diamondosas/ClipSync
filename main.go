package main

import (

	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"

	
	"clipsync/internal/clipboard"
	"clipsync/internal/globals"
	"clipsync/internal/network"
	"golang.org/x/sync/errgroup"
	// "clipsync/internal/ping"
)

var Version = "dev"

func main() {

	clipboard.Init()
	network.RegisterDevice()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	globals.WG.Add(4)
	fmt.Println(1)
	go RunWithContext(ctx, network.BrowseForDevices)
	fmt.Println(2)
	go RunWithContext(ctx, network.Listen)
	fmt.Println(3)
	go func() {
		defer globals.WG.Done()
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
		defer globals.WG.Done()
		<-network.Ready
		for {
			buffer, n := network.RecieveClipboard()
			clipboard.WriteClipboard(string(buffer[:n]))
		}
	}()

	globals.WG.Wait()
}

func RunWithContext(ctx context.Context, task func()) {
	defer globals.WG.Done()
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
