package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"clipsync/internal/clipboard"
	"clipsync/internal/globals"
	"clipsync/internal/network"
	// "clipsync/internal/ping"
)

func main() {
	
	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	globals.WG.Add(5)
	go RunWithContext(ctx, network.RegisterDevice)
	fmt.Println(1)
	go RunWithContext(ctx, network.BrowseForDevices)
	fmt.Println(2)
	go RunWithContext(ctx, network.Listen)
	fmt.Println(3)
	go func(){
        changedText := clipboard.ChangedClipboard(ctx) // make this return the channel
        for data := range changedText {
            network.SendData(data)
		}
    }()
	go func(){
		<-network.Ready
		buffer, n := network.RecieveClipboard()
		clipboard.WriteClipboard(string(buffer[:n]))
	}()
	
	globals.WG.Wait()	
}


func RunWithContext(ctx context.Context, task func()){
	defer globals.WG.Done()
	go task()
	<-ctx.Done()
}