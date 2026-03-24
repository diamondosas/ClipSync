package view

import(
	"log"
	"clipsync/internal/globals"
	"clipsync/gui"
)
// It shooud be a goroutune and it should work in such a way that the fucntion is kickstaed bt a pointer 
// It should then like update the vie which then redrws the GUI tell me if my implemtation is wrong 
func UpdateDevices(Device globals.Device){
	// It hssould change the 
	globals.ConnDevices = append(globals.ConnDevices, Device)
	//Make this equal to the Devices : []pages.Deivce so that i t can change it 
	RedrawUI()
}

func UpdateClipborad(data string){
	
	
	RedrawUI()
}

func RedrawUI(){
	//Redraw the UI to sho wchanges in boht Update Deivecea dn clipBoard
}