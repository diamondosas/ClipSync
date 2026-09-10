package gui

import (
	"clipsync/gui/pages"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)
var State *AppState
// AppState keeps track of the global application state.
// This struct ensures our GUI is interactive and holds the mock data.
type AppState struct {
	Theme *material.Theme

	// Navigation
	ActiveTab int
	TabBtns   [2]widget.Clickable

	// Devices Page State
	DeviceList widget.List
	Devices    []pages.Device

	// Clipboard Page State
	ClipList widget.List
	History  []string

	// Dialog State
	HelpBtn      widget.Clickable
	CloseHelpBtn widget.Clickable
	ShowHelp     bool

	// Channel For Thread-Safe UI Updates
	DeviceUpdates      chan pages.Device
	DevicePruneUpdates chan []string
	ClipUpdates        chan string
}

// NewAppState initializes the default state of the App.
func NewAppState(th *material.Theme) *AppState {
	s := &AppState{
		Theme:              th,
		Devices:            []pages.Device{},
		History:            []string{},
		DeviceUpdates:      make(chan pages.Device, 50),
		DevicePruneUpdates: make(chan []string, 10),
		ClipUpdates:        make(chan string, 50),
	}
	// Setup Lists to be Vertical
	s.DeviceList.Axis = layout.Vertical
	s.ClipList.Axis = layout.Vertical
	State = s
	return s
}

func (s *AppState) UpdateUIValues() {
	// 1. Drain new device updates
	for {
		select {
		case dev := <-s.DeviceUpdates:
			exists := false
			for _, d := range s.Devices {
				if d.IP == dev.IP {
					exists = true
					break
				}
			}
			if !exists {
				s.Devices = append(s.Devices, dev)
			}
		default:
			goto CheckPrune
		}
	}

CheckPrune:
	// 2. Drain device pruning updates
	for {
		select {
		case activeIPs := <-s.DevicePruneUpdates:
			activeMap := make(map[string]bool, len(activeIPs))
			for _, ip := range activeIPs {
				activeMap[ip] = true
			}
			var remaining []pages.Device
			for _, d := range s.Devices {
				if activeMap[d.IP] {
					remaining = append(remaining, d)
				}
			}
			s.Devices = remaining
		default:
			goto CheckClip
		}
	}

CheckClip:
	// 3. Drain clipboard history updates
	for {
		select {
		case clip := <-s.ClipUpdates:
			s.History = append([]string{clip}, s.History...)
		default:
			return
		}
	}
}

// Update processes any events/clicks before layout rendering.
func (s *AppState) Update(gtx layout.Context) {
	
	s.UpdateUIValues()

	// Handle Tab Clicks
	for i := range s.TabBtns {
		if s.TabBtns[i].Clicked(gtx) {
			s.ActiveTab = i
		}
	}

	// Handle Help Dialog
	if s.HelpBtn.Clicked(gtx) {
		s.ShowHelp = true
	}
	if s.CloseHelpBtn.Clicked(gtx) {
		s.ShowHelp = false
	}
}
