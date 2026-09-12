package tray

import (
	_ "embed"
	"log"
	"os"
	"context"
	"runtime"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var trayIcon_ico []byte

//go:embed icon.png
var trayIcon_png []byte
func OnTrayReady(cancel context.CancelFunc,onOpenWindow func()){

	if runtime.GOOS == "windows"{
		systray.SetIcon(trayIcon_ico)
	}else{
		systray.SetIcon(trayIcon_png)
	}
	systray.SetTitle("Clipsync")
	systray.SetTooltip("Clipsync is Running")

	open := systray.AddMenuItem("Open ", "Open Clipsync")
	quit := systray.AddMenuItem("Quit", "Quit Clipsync")

	go func(){
		for{
			select{
			case <-open.ClickedCh:
				if onOpenWindow != nil{
					onOpenWindow()
				}
			case <-quit.ClickedCh:
				cancel()
				systray.Quit() 
				return
			}
		}

	}()
}

func OnTrayExit(){
	log.Println("Exiting")
	os.Exit(0)
}