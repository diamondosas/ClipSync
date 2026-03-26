	package view

import (
	"clipsync/gui"
	"clipsync/gui/pages"
	"clipsync/internal/globals"

)



// It shooud be a goroutune and it should work in such a way that the fucntion is kickstaed bt a pointer 
// It should then like update the vie which then redrws the GUI tell me if my implemtation is wrong 
func UpdateDevices(Device globals.Device) {
	// 1. Create a GUI device from the backend device
	newDevice := pages.Device{
		Name: Device.Name,
		IP:   Device.Ip,
	}

	// 2. Add it to our State (using the global pointer State)
	gui.State.Devices = append(gui.State.Devices, newDevice)

	// 3. Tell the UI to redraw
	RedrawUI()
}

func UpdateClipboard(data string) {
	// It is not meant to be append because i want it to act like a stack
	// This adds the new item to the FRONT of the list
	if data != ""{
		gui.State.History = append([]string{data}, gui.State.History...)

		RedrawUI()
	}
}

func RedrawUI() {
	// Redraw the UI to show changes in both Update Devices and Clipboard
	if gui.Window != nil{
		gui.Window.Invalidate()
	}
}