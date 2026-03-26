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
}

// NewAppState initializes the default state of the App.
func NewAppState(th *material.Theme) *AppState {
	s := &AppState{
		Theme: th,
		Devices: []pages.Device{
			// {Name: "Desktop-PC", IP: "192.168.1.10"},
			// {Name: "MacBook-Pro", IP: "192.168.1.12"},
			// {Name: "Android-Phone", IP: "192.168.1.15"},
		},
		History: []string{
			// "Hello World!",
			// "https://github.com/leojimenezg/scapmi",
			// "func main() { fmt.Println(GUI Rocks) }",
			// "Mock Clipboard Data 4",
			// "Mock Clipboard Data 5",
			// "Mock Clipboard Data 6",
			// "Mock Clipboard Data 7",
		},
	}
	// Setup Lists to be Vertical
	s.DeviceList.Axis = layout.Vertical
	s.ClipList.Axis = layout.Vertical
	State = s
	return s
}

// Update processes any events/clicks before layout rendering.
func (s *AppState) Update(gtx layout.Context) {
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
