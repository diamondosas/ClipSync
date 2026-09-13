package view

import (
	"clipsync/gui"
	"clipsync/gui/pages"
	"clipsync/internal"
)

// AddNewDevice sends a new device to the global store and the GUI channel
func AddNewDevice(device internal.Device) {
	// 1. Update Global State
	internal.ConnDevicesMu.Lock()
	for i := range internal.ConnDevices {
		if internal.ConnDevices[i] == device {
			return
		}
	}
	internal.ConnDevices = append(internal.ConnDevices, device)
	internal.ConnDevicesMu.Unlock()

	// 2. Send to GUI channel if GUI is running
	if gui.State != nil && gui.State.DeviceUpdates != nil {
		select {
		case gui.State.DeviceUpdates <- pages.Device{Name: device.Name, IP: device.Ip}:
			RedrawUI()
		default:
		}
	}
}

// UpdateClipboard sends new clipboard text to the global history and the GUI channel
func UpdateClipboard(data string) {
	if data == "" {
		return
	}

	// 1. Update Global State
	internal.ClipHistoryMu.Lock()
	internal.ClipHistory = append([]string{data}, internal.ClipHistory...)
	internal.ClipHistoryMu.Unlock()

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

// UpdateClipboardSynced sends incoming clipboard text from the network and triggers a sync notification
func UpdateClipboardSynced(data string) {
	UpdateClipboard(data)
	if gui.State != nil && gui.State.ToastUpdates != nil {
		select {
		case gui.State.ToastUpdates <- "Synced new clip":
			RedrawUI()
		default:
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
	internal.ConnDevicesMu.Lock()
	var updatedGlobal []internal.Device
	for _, d := range internal.ConnDevices {
		if activeMap[d.Ip] {
			updatedGlobal = append(updatedGlobal, d)
		}
	}
	internal.ConnDevices = updatedGlobal
	internal.ConnDevicesMu.Unlock()

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
