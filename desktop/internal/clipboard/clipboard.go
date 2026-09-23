package clipboard

import (
	"context"
	"log"
	"slices"

	"clipsync/internal"

	"golang.design/x/clipboard"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func CopyClipboard(ctx context.Context) string {
	data, err := clipboard.Read(ctx, clipboard.FmtText)
	if err != nil {
		log.Println("Could not read Clipboard")
	}
	return string(data)
}

func WriteClipboard(ctx context.Context, data string) {
	buf := []byte(data)
	_, _ = clipboard.Write(ctx, clipboard.FmtText, buf)
}

func WatchClipboard(ctx context.Context) <-chan []byte {
	text := clipboard.Watch(ctx, clipboard.FmtText)
	var out = make(chan []byte, 1)
	
	//REFACTOR: Put in root.go so that it ca avoid it importing network module
	go func() {
		for {
			select {
			case data := <-text:
				internal.LastRecvClipMu.Lock()
				isEqual := slices.Equal(data.Bytes, internal.LastRecvClip)
				internal.LastRecvClipMu.Unlock()
				if !isEqual {
					select {
					case out <- data.Bytes:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
