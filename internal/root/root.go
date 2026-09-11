package root

import (
	"context"
	"log"
	"time"

	"clipsync/internal/clipboard"
	"clipsync/internal/globals"
	"clipsync/internal/network"
	"clipsync/internal/view"

	"golang.org/x/sync/errgroup"
)

// initializes and runs all background synchronization tasks.
func StartClipSync(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	// 1. Register our device on the network
	eg.Go(func() error {
		return network.RegisterDevice(ctx, "")
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
		var clipCh = make(chan string, len(clip))
		if clip == nil {
			log.Println("[Sync] Warning: clipboard watch channel is nil (clipboard may not be initialized)")
			<-ctx.Done()
			return ctx.Err()
		}

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case data, ok := <-clipCh:
				if !ok {
					return nil
				}
				if len(data) == 0 {
					continue
				}
				// Avoid loops: don't send if it's the same as what we just received
					log.Printf("[Sync] Local change detected, sending to %d devices", len(globals.IPS))
					network.SendClipboard([]byte(data))
					view.UpdateClipboard(string(data))
			}
		}
	})

	// 5. Receive clipboard data from other devices
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
				buffer, n := network.RecieveClipboard()
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
			case <-time.After(5 * time.Second):
				
				globals.IPSMu.Lock()
				var ipsToPing []string
				copy(ipsToPing, globals.IPS)
				globals.IPSMu.Unlock()

				if len(ipsToPing) > 0 {
					currentIPS := network.PingIPS(ipsToPing)

					globals.IPSMu.Lock()
					globals.IPS = currentIPS
					globals.IPSMu.Unlock()
					view.PruneDevices(currentIPS)
					log.Printf("[Sync] Ping check: %d devices active", len(globals.IPS))
				}
			}
		}
	})

	return eg.Wait()
}
