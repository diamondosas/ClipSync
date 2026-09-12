// ai-generated
package gui

import (
	"context"
	"log"

	"clipsync/gui/pages"
	"clipsync/gui/utils"
	"golang.design/x/clipboard"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var State *AppState

// AppState keeps track of the global application state.
type AppState struct {
	Theme *material.Theme

	// Navigation
	ActiveTab int
	TabBtns   [2]widget.Clickable

	// Devices Page State
	DeviceList widget.List
	Devices    []pages.Device

	// Clipboard Page State
	ClipList      widget.List
	ClipItems     []*pages.ClipItem
	ClearClipsBtn widget.Clickable

	// Dialog State
	HelpBtn      widget.Clickable
	CloseHelpBtn widget.Clickable
	ShowHelp     bool

	// Channels For Thread-Safe UI Updates
	DeviceUpdates      chan pages.Device
	DevicePruneUpdates chan []string
	ClipUpdates        chan string
}

// NewAppState initializes the default state of the App.
func NewAppState(th *material.Theme) *AppState {
	s := &AppState{
		Theme:              th,
		Devices:            []pages.Device{},
		ClipItems:          []*pages.ClipItem{},
		DeviceUpdates:      make(chan pages.Device, 50),
		DevicePruneUpdates: make(chan []string, 10),
		ClipUpdates:        make(chan string, 50),
	}
	// Setup Lists to be Vertical
	s.DeviceList.Axis = layout.Vertical
	s.ClipList.Axis = layout.Vertical

	// Load pinned clips from persistent disk storage
	for _, text := range utils.LoadPinned() {
		s.ClipItems = append(s.ClipItems, &pages.ClipItem{
			Content:  text,
			IsPinned: true,
		})
	}

	State = s
	return s
}

// UpdateUIValues drains background channels into UI slices.
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
	// 3. Drain clipboard updates
	for {
		select {
		case clip := <-s.ClipUpdates:
			s.AddClip(clip)
		default:
			return
		}
	}
}

// AddClip adds an incoming clipboard item below pinned items and deduplicates.
func (s *AppState) AddClip(text string) {
	if text == "" {
		return
	}
	// Do not duplicate if already pinned
	for _, item := range s.ClipItems {
		if item.IsPinned && item.Content == text {
			return
		}
	}
	// Remove previous unpinned instance if it exists to bring to top
	var filtered []*pages.ClipItem
	for _, item := range s.ClipItems {
		if !item.IsPinned && item.Content == text {
			continue
		}
		filtered = append(filtered, item)
	}
	s.ClipItems = filtered

	newItem := &pages.ClipItem{
		Content:  text,
		IsPinned: false,
	}

	// Insert right after pinned items
	insertIdx := 0
	for _, item := range s.ClipItems {
		if item.IsPinned {
			insertIdx++
		} else {
			break
		}
	}
	s.ClipItems = append(s.ClipItems[:insertIdx], append([]*pages.ClipItem{newItem}, s.ClipItems[insertIdx:]...)...)
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

	// Handle Clear All Button
	if s.ClearClipsBtn.Clicked(gtx) {
		s.ClearUnpinnedClips()
	}

	// Handle Clipboard Item Actions (Copy, Pin, Delete)
	s.HandleClipEvents(gtx)
}

// HandleClipEvents checks clicks on individual clipboard cards.
func (s *AppState) HandleClipEvents(gtx layout.Context) {
	var remaining []*pages.ClipItem
	pinnedChanged := false

	for _, item := range s.ClipItems {
		// 1. Delete clicked
		if item.DeleteBtn.Clicked(gtx) {
			if item.IsPinned {
				pinnedChanged = true
			}
			continue
		}

		// 2. Pin toggled
		if item.PinBtn.Clicked(gtx) {
			item.IsPinned = !item.IsPinned
			pinnedChanged = true
		}

		// 3. Card clicked -> write to system clipboard
		if item.CardBtn.Clicked(gtx) {
			_, _ = clipboard.Write(context.Background(), clipboard.FmtText, []byte(item.Content))
		}

		remaining = append(remaining, item)
	}

	s.ClipItems = remaining

	if pinnedChanged {
		s.sortClips()
		s.persistPinned()
	}
}

// sortClips ensures all pinned items remain at the top of the list.
func (s *AppState) sortClips() {
	var pinned []*pages.ClipItem
	var unpinned []*pages.ClipItem
	for _, item := range s.ClipItems {
		if item.IsPinned {
			pinned = append(pinned, item)
		} else {
			unpinned = append(unpinned, item)
		}
	}
	s.ClipItems = append(pinned, unpinned...)
}

// ClearUnpinnedClips removes all non-pinned clips from memory.
func (s *AppState) ClearUnpinnedClips() {
	var pinned []*pages.ClipItem
	for _, item := range s.ClipItems {
		if item.IsPinned {
			pinned = append(pinned, item)
		}
	}
	s.ClipItems = pinned
}

// persistPinned writes the current pinned items to disk.
func (s *AppState) persistPinned() {
	var texts []string
	for _, item := range s.ClipItems {
		if item.IsPinned {
			texts = append(texts, item.Content)
		}
	}
	go func(saved []string) {
		if err := utils.SavePinned(saved); err != nil {
			log.Printf("[Storage] Error saving pinned clips: %v", err)
		}
	}(texts)
}
