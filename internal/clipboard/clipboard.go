package clipboard

import (
	"context"
	"log"

	// "sync"
	"golang.design/x/clipboard"
	// "clipsync/internal/globals"
)

func init() {
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

func ChangedClipbord(ctx context.Context){
	changedText := clipboard.Watch(ctx, clipboard.FmtText)
	for data := range changedText{
		WriteClipboard(string(data))
	}
}
