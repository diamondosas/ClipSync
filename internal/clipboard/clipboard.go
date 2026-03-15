package clipboard

import (
	"context"
	"log"

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

func ChangedClipboard(ctx context.Context) <-chan []byte{
	changedText := clipboard.Watch(ctx, clipboard.FmtText)
	return changedText
}
