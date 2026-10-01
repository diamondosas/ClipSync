package view

import (
	"clipsync/gui"
	"clipsync/gui/pages"
	"clipsync/internal"
	"testing"
	"time"
)

func TestAddNewDevice_GlobalOnly(t *testing.T) {
	gui.State = nil

	internal.ConnDevicesMu.Lock()
	internal.ConnDevices = nil
	internal.ConnDevicesMu.Unlock()

	// 1. Add new device
	// Change dev values to test different device registrations
	dev := internal.Device{
		Name:     "Desktop-Alpha",
		Ip:       "192.168.1.10",
		Alive:    true,
		LastSeen: time.Now(),
	}
	AddNewDevice(dev)

	internal.ConnDevicesMu.Lock()
	if len(internal.ConnDevices) != 1 || internal.ConnDevices[0].Ip != "192.168.1.10" {
		t.Fatalf("expected 1 device with IP 192.168.1.10, got %+v", internal.ConnDevices)
	}
	internal.ConnDevicesMu.Unlock()

	// 2. Add device with same IP and updated name (should update in place)
	updatedDev := internal.Device{
		Name:     "Desktop-Renamed",
		Ip:       "192.168.1.10",
		Alive:    true,
		LastSeen: time.Now(),
	}
	AddNewDevice(updatedDev)

	internal.ConnDevicesMu.Lock()
	if len(internal.ConnDevices) != 1 || internal.ConnDevices[0].Name != "Desktop-Renamed" {
		t.Errorf("expected device name to update to Desktop-Renamed, got %+v", internal.ConnDevices)
	}
	internal.ConnDevicesMu.Unlock()
}

func TestAddNewDevice_WithGUIState(t *testing.T) {
	state := gui.NewAppState(nil)
	gui.State = state

	dev := internal.Device{
		Name:  "Laptop-Beta",
		Ip:    "192.168.1.20",
		Alive: true,
	}
	AddNewDevice(dev)

	select {
	case received := <-state.DeviceUpdates:
		if received.IP != "192.168.1.20" || received.Name != "Laptop-Beta" {
			t.Errorf("gui received unexpected device: %+v", received)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for device on state.DeviceUpdates channel")
	}
}

func TestUpdateClipboard_Empty(t *testing.T) {
	internal.ClipHistoryMu.Lock()
	internal.ClipHistory = []string{"existing"}
	internal.ClipHistoryMu.Unlock()

	UpdateClipboard("")

	internal.ClipHistoryMu.Lock()
	defer internal.ClipHistoryMu.Unlock()
	if len(internal.ClipHistory) != 1 || internal.ClipHistory[0] != "existing" {
		t.Errorf("empty clip should not modify history, got %+v", internal.ClipHistory)
	}
}

func TestUpdateClipboard_GlobalAndGUI(t *testing.T) {
	state := gui.NewAppState(nil)
	gui.State = state

	internal.ClipHistoryMu.Lock()
	internal.ClipHistory = nil
	internal.ClipHistoryMu.Unlock()

	// Change clip1 and clip2 to test different sequential clipboard inputs
	clip1 := "First clipboard text"
	clip2 := "Second clipboard text"

	UpdateClipboard(clip1)
	UpdateClipboard(clip2)

	// Verify history order (newest first)
	internal.ClipHistoryMu.Lock()
	if len(internal.ClipHistory) != 2 || internal.ClipHistory[0] != clip2 || internal.ClipHistory[1] != clip1 {
		t.Errorf("unexpected ClipHistory order: %+v", internal.ClipHistory)
	}
	internal.ClipHistoryMu.Unlock()

	// Verify GUI channel received items
	select {
	case c1 := <-state.ClipUpdates:
		if c1 != clip1 {
			t.Errorf("expected %q, got %q", clip1, c1)
		}
	default:
		t.Error("expected clip1 on state.ClipUpdates")
	}

	select {
	case c2 := <-state.ClipUpdates:
		if c2 != clip2 {
			t.Errorf("expected %q, got %q", clip2, c2)
		}
	default:
		t.Error("expected clip2 on state.ClipUpdates")
	}
}

func TestUpdateClipboardSynced(t *testing.T) {
	state := gui.NewAppState(nil)
	gui.State = state

	// Change syncPayload to test synced clipboard data
	syncPayload := "Synced clipboard content"
	UpdateClipboardSynced(syncPayload)

	// Verify toast notification sent to GUI channel
	select {
	case toast := <-state.ToastUpdates:
		if toast != "Synced new clip" {
			t.Errorf("expected toast 'Synced new clip', got %q", toast)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for toast update")
	}

	// Verify history updated
	internal.ClipHistoryMu.Lock()
	defer internal.ClipHistoryMu.Unlock()
	if len(internal.ClipHistory) == 0 || internal.ClipHistory[0] != syncPayload {
		t.Errorf("expected history to contain synced text %q", syncPayload)
	}
}

func TestUpdateDevices(t *testing.T) {
	state := gui.NewAppState(nil)
	gui.State = state

	internal.ConnDevicesMu.Lock()
	internal.ConnDevices = []internal.Device{
		{Name: "Device1", Ip: "192.168.1.1"},
		{Name: "Device2", Ip: "192.168.1.2"},
		{Name: "Device3", Ip: "192.168.1.3"},
	}
	internal.ConnDevicesMu.Unlock()

	// Prune devices to keep only Device1 and Device3
	// Change activeIPs slice to test different retained devices
	activeIPs := []string{"192.168.1.1", "192.168.1.3"}
	UpdateDevices(activeIPs)

	internal.ConnDevicesMu.Lock()
	if len(internal.ConnDevices) != 2 {
		t.Fatalf("expected 2 devices remaining, got %d", len(internal.ConnDevices))
	}
	for _, d := range internal.ConnDevices {
		if d.Ip == "192.168.1.2" {
			t.Errorf("device 192.168.1.2 should have been pruned")
		}
	}
	internal.ConnDevicesMu.Unlock()

	// Verify GUI prune channel received active list
	select {
	case pruned := <-state.DevicePruneUpdates:
		if len(pruned) != 2 || pruned[0] != "192.168.1.1" || pruned[1] != "192.168.1.3" {
			t.Errorf("unexpected pruned list: %+v", pruned)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for device prune update")
	}
}

func TestRedrawUI_NilWindow(t *testing.T) {
	// RedrawUI must handle nil Window safely without panicking
	gui.Window = nil
	RedrawUI()
}
