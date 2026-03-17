package clipboard

import (
	"clipsync/internal/network"
	"context"
	"log"
	"slices"
	// "sync"
	"golang.design/x/clipboard"
	// "clipsync/internal/network"
)

func Init() {	
	err := clipboard.Init()
	if err != nil {
		log.Println(err)
	}
}

func CopyClipboard() string {
	data := clipboard.Read(clipboard.FmtText)
	return string(data)
}

func WriteClipboard(data string) {
	byte := []byte(data)
	clipboard.Write(clipboard.FmtText, byte)
}

func WatchClipboard(ctx context.Context, changedText chan<- []byte){
	text := clipboard.Watch(ctx, clipboard.FmtText)
	for{
		select{
		case data := <-text:
			if !slices.Equal(data, network.Buffer){
				changedText <- data
			}
		case <-ctx.Done():
			return
		}
	}
}
