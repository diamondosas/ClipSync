package view

import (
	"clipsync/gui"
	"clipsync/internal/globals"
	"log"

	"gioui.org/app"
)

var State *gui.AppState
var Window *app.Window
// It shooud be a goroutune and it should work in such a way that the fucntion is kickstaed bt a pointer 
// It should then like update the vie which then redrws the GUI tell me if my implemtation is wrong 
func UpdateDevices(Device globals.Device){
	// It hssould change the 
	globals.ConnDevices = append(globals.ConnDevices, Device)
	//Make this equal to the Devices : []pages.Deivce so that i t can change it 
	RedrawUI()
	globals.ConnDevices = gui.AppState.Devices
}

func UpdateClipborad(data string){
	//It is not menat to be append beacuse i wnat it to act like a stack or maybe the gui will just reverse how it is represented
	gui.AppState.History = append(gui.AppState.History, data)
	RedrawUI()
}

func RedrawUI(){
	//Redraw the UI to sho wchanges in boht Update Deivecea dn clipBoard
}