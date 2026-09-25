package root

import (
	"bytes"
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

	// 3. Watch local clipboard for changes
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
				internal.LastRecvClipMu.Lock()
				isEcho := bytes.Equal(data, internal.LastRecvClip)
				internal.LastRecvClipMu.Unlock()
				if isEcho {
					continue
				}

				internal.ConnDevicesMu.Lock()
				for i := range internal.ConnDevices {
					if internal.ConnDevices[i].Alive {
						log.Println("[Sync] Local change detected, sending to Device: ", internal.ConnDevices[i].Ip)
					}
				}
				internal.ConnDevicesMu.Unlock()

				network.SendClipboard([]byte(data))
				view.UpdateClipboard(string(data))
			}
		}
	})

	// 4. Receive Clipboard data from peers
	eg.Go(func() error {
		clips := network.ReceiveClipboard(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case clip, ok := <-clips:
				if !ok {
					return nil
				}
				if len(clip) > 0 {
					data := string(clip)
					log.Printf("[Sync] Received new clipboard data (%d bytes)", len(clip))
					clipboard.WriteClipboard(ctx, data)
					view.UpdateClipboardSynced(data)
				}
			}
		}
	})

	// 5. Periodically ping devices to tell them I am still alive
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
				internal.ConnDevicesMu.Lock()
				var ipsToPing []string
				for i := range internal.ConnDevices {
					ipsToPing = append(ipsToPing, internal.ConnDevices[i].Ip)
				}
				internal.ConnDevicesMu.Unlock()

				if len(ipsToPing) > 0 {
					network.PingIPS(ipsToPing)
				}
			}
		}
	})

	// 6. Periodically check whether devices in the ConnDevices list are still alive
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Second):
				internal.ConnDevicesMu.Lock()
				hasDevices := len(internal.ConnDevices) > 0
				internal.ConnDevicesMu.Unlock()

				if hasDevices {
					network.CheckForPing()
				}
			}
		}
	})

	return eg.Wait()
}
