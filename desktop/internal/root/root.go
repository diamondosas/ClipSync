package root

import (
	"context"
	"log"
	"time"

	"clipsync/internal"
	"clipsync/internal/clipboard"
	"clipsync/internal/network"
	"clipsync/internal/view"

	"golang.org/x/sync/errgroup"
)

// initializes and runs all background synchronization tasks.
func StartClipSync(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	// 1. Auto-discover and register dynamically
	eg.Go(func() error {
		return network.StartAutoDiscovery(ctx)
	})

	// 2. Listen for incoming UDP connections
	eg.Go(func() error {
		return network.Listen(ctx)
	})

	// 4. Watch local clipboard for changes
	eg.Go(func() error {
		clip := clipboard.WatchClipboard(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case data, ok := <-clip:
				if !ok {
					return nil
				}
				if len(data) == 0 {
					continue
				}
				// Avoid loops: don't send if it's the same as what we just received
				internal.ConnDevicesMu.Lock()
				for i := range internal.ConnDevices{
					log.Println("[Sync] Local change detected, sending to Device: ", internal.ConnDevices[i].Ip)
				}
				internal.ConnDevicesMu.Unlock()
				network.SendClipboard([]byte(data))
				view.UpdateClipboard(string(data))
			}
		}
	})

	// 5. Receive Clipboard data from each peers on the network (it reieves from session independently)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case clip := <- network.ReceiveClipboard(ctx):
				if len(clip) > 0 {
					data := string(clip)
					log.Printf("[Sync] Received new clipboard data (%d bytes)", len(clip))
					clipboard.WriteClipboard(ctx, data)
					view.UpdateClipboardSynced(data)
				}
			}
		}
	})

	// 6. Periodically ping devices to tell them I am still alive
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
				internal.ConnDevicesMu.Lock()
				var ipsToPing []string
				if internal.ConnDevices != nil{
					for i:= range internal.ConnDevices{
						ipsToPing = append(ipsToPing, internal.ConnDevices[i].Ip)
					}
				}

				internal.ConnDevicesMu.Unlock()

				if len(ipsToPing) > 0 {
					network.PingIPS(ipsToPing)
				}
			}
		}
	})

	// Perodically Check whether devices in the []ConnDevice list has sent that they are alive
	eg.Go(func() error{
		for{
			select{
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Second):
				internal.ConnDevicesMu.Lock()
				if internal.ConnDevices != nil{
					internal.ConnDevicesMu.Unlock()
					log.Println("Checking for Dead Connection")
					network.CheckForPing()
				}
			}
		}
	} )

	return eg.Wait()
}
