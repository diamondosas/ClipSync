package tray

import (
	_ "embed"
	"log"
	"os"
	"context"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var trayIcon []byte

func OnTrayReady(cancel context.CancelFunc){
	systray.SetIcon(trayIcon)
	systray.SetTitle("Clipsync")
	systray.SetTooltip("Clipsync is Running")


	quit := systray.AddMenuItem("Quit", "Quit Clipsync")

	go func(){
		<-quit.ClickedCh
		cancel()
		systray.Quit()
	}()
}

func OnTrayExit(){
	log.Println("Exiting")
	os.Exit(0)
}