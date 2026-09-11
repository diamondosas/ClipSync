package clipboard

import (
	"context"
	"log"

	"clipsync/internal/network"
	"golang.design/x/clipboard"
)

func init()  {	
	err := clipboard.Init()
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func CopyClipboard(ctx context.Context) string {
	data, err  := clipboard.Read(ctx, clipboard.FmtText)
	if err != nil{
		log.Println("Could not read Clipboard")
	}
	return string(data)
}

func WriteClipboard(ctx context.Context, data string) {
	byte := []byte(data)
	_, _ = clipboard.Write(ctx, clipboard.FmtText, byte)
}

func WatchClipboard(ctx context.Context) <-chan []byte {
	text := clipboard.Watch(ctx, clipboard.FmtText)
	var out = make(chan []byte, 1)

	go func(){
		for {
			select {
			case data:= <-text:
				if !network.IsLastReceived(data.Bytes){
					select{
					case out <- data.Bytes:
					}
				}
			case <-ctx.Done():
				return 
			}		
		}
	}()

	return out	
}