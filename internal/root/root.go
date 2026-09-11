package root

import (
	"context"
	"log"
	"time"

	"clipsync/internal/clipboard"
	"clipsync/internal/network"
	"clipsync/internal/view"

	"golang.org/x/sync/errgroup"
)

// initializes and runs all background synchronization tasks.
func StartClipSync(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	// 1. Register our device on the network
	eg.Go(func() error {
		return network.RegisterDevice(ctx)
	})

	// 2. Discover other devices
	eg.Go(func() error {
		return network.BrowseForDevices(ctx)
	})

	// 3. Listen for incoming UDP connections
	eg.Go(func() error {
		return network.Listen(ctx)
	})

	// 4. Watch local clipboard for changes
	eg.Go(func() error {
		clip := clipboard.WatchClipboard(ctx)
		if clip == nil {
			log.Println("[Sync] Warning: clipboard watch channel is nil (clipboard may not be initialized)")
			<-ctx.Done()
			return ctx.Err()
		}

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
				log.Printf("[Sync] Local change detected, sending to %d devices", len(IPS))
				network.SendClipboard([]byte(data))
				view.UpdateClipboard(string(data))
			}
		}
	})

	// 5. Receive data, clipboard and pings from other devices and Update Clipboard if it contain the Data
	eg.Go(func() error {
		select {
		case <-network.Ready:
		case <-ctx.Done():
			return ctx.Err()
		}
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				//If the Data is Clipboard it iwll
				buffer, n := network.RecieveData()

				if n > 0 {
					data := string(buffer[:n])
					log.Printf("[Sync] Received new clipboard data (%d bytes)", n)
					clipboard.WriteClipboard(ctx, data)
					view.UpdateClipboard(data)
				}
			}
		}
	})

	// 6. Periodically ping devices to keep the list fresh
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Second):

				network.IPSMu.Lock()
				var ipsToPing []string
				copy(ipsToPing, network.IPS)
				network.IPSMu.Unlock()

				if len(ipsToPing) > 0 {
					network.PingIPS(ipsToPing)
				}
			}
		}
	})

	return eg.Wait()
}
