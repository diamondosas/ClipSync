package view

import (
	"clipsync/gui"
	"clipsync/gui/pages"
	"clipsync/internal/globals"
)

// AddNewDevice sends a new device to the global store and the GUI channel
func AddNewDevice(device globals.Device) {
	// 1. Update Global State
	globals.ConnDevicesMu.Lock()
	globals.ConnDevices = append(globals.ConnDevices, device)
	globals.ConnDevicesMu.Unlock()

	// 2. Send to GUI channel if GUI is running
	if gui.State != nil && gui.State.DeviceUpdates != nil {
		select {
		case gui.State.DeviceUpdates <- pages.Device{Name: device.Name, IP: device.Ip}:
			RedrawUI()
		}
	}
}

// UpdateClipboard sends new clipboard text to the global history and the GUI channel
func UpdateClipboard(data string) {
	if data == "" {
		return
	}

	// 1. Update Global State
	globals.ClipHistoryMu.Lock()
	globals.ClipHistory = append([]string{data}, globals.ClipHistory...)
	globals.ClipHistoryMu.Unlock()

	// 2. Send to GUI channel if GUI is running
	if gui.State != nil && gui.State.ClipUpdates != nil {
		select {
		case gui.State.ClipUpdates <- data:
			RedrawUI()
		default:
			// Channel buffer full; drop or log to prevent blocking
		}
	}
}

// UpdateDevices updates the global and GUI device list based on reachable IPs
func UpdateDevices(activeIPs []string) {
	activeMap := make(map[string]bool, len(activeIPs))
	for _, ip := range activeIPs {
		activeMap[ip] = true
	}

	// 1. Update Global State
	globals.ConnDevicesMu.Lock()
	var updatedGlobal []globals.Device
	for _, d := range globals.ConnDevices {
		if activeMap[d.Ip] {
			updatedGlobal = append(updatedGlobal, d)
		}
	}
	globals.ConnDevices = updatedGlobal
	globals.ConnDevicesMu.Unlock()

	// 2. Send pruned device list to GUI channel
	if gui.State != nil && gui.State.DevicePruneUpdates != nil {
		select {
		case gui.State.DevicePruneUpdates <- activeIPs:
			RedrawUI()
		default:
		}
	}
}

func RedrawUI() {
	if gui.Window != nil {
		gui.Window.Invalidate()
	}
}
