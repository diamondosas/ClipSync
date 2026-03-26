package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"
	
	"clipsync/gui"
	"clipsync/internal/clipboard"
	"clipsync/internal/globals"
	"clipsync/internal/network"
	"clipsync/internal/ping"
	"clipsync/internal/view"

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
		return network.Listen(ctx)
	})

	eg.Go(func() error{
		for {
			data := clipboard.WatchClipboard(ctx)
			if !slices.Equal(data, network.Buffer) {
				network.SendClipboard(data)
				view.UpdateClipborad(string(data))
			}
		}
	})
	eg.Go(func() error {
		<-network.Ready
		for {
			buffer, n := network.RecieveClipboard()
			data := string(buffer[:n])
			clipboard.WriteClipboard(data)
			view.UpdateClipborad(data)
		}
	})
	
	eg.Go(func() error {
		for {
			globals.IPSMu.Lock()
			ipsToPing := make([]string, len(globals.IPS))
			copy(ipsToPing, globals.IPS)
			globals.IPSMu.Unlock()
			
			currentIPS := ping.PingIPS(ipsToPing)

			globals.IPSMu.Lock()
			globals.IPS = currentIPS
			globals.IPSMu.Unlock()
		}

	})

	eg.Go(func() error{
		err := eg.Wait()
		if err != nil {
			log.Println(err)
		}
		return err
	})
	
	gui.StartGUI()
	
}
