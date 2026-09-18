package gui

import (
	"clipsync-android/gui/pages"
	"clipsync-android/gui/utils"
	"clipsync-android/internal/events"
	"clipsync-android/internal/service"
	"io"
	"strings"
	"sync"
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/transfer"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type AppState struct {
	Theme        *material.Theme
	Service      *service.ClipSyncService
	InvalidateFn func()

	// Navigation
	ActiveTab int
	TabBtns   [2]widget.Clickable

	// Devices State
	DeviceList widget.List
	DevicesMu  sync.Mutex
	Devices    []pages.DeviceViewItem

	// Clipboard State
	ClipList       widget.List
	ClipMu         sync.Mutex
	ClipItems      []*pages.ClipItem
	ClearClipsBtn  widget.Clickable
	SyncClipBtn    widget.Clickable
	SearchEditor   widget.Editor
	SearchClearBtn widget.Clickable

	// Toast State
	ToastMsg       string
	ToastStartTime time.Time

	// Dialog State
	HelpBtn            widget.Clickable
	CloseHelpBtn       widget.Clickable
	ShowHelp           bool
	ConnectBtn         widget.Clickable
	DialogConnectBtn   widget.Clickable
	DialogCancelBtn    widget.Clickable
	IPEditor           widget.Editor
	ShowConnectDialog  bool

	subs []events.Subscription
}

func NewAppState(th *material.Theme, svc *service.ClipSyncService, invalidateFn func()) *AppState {
	s := &AppState{
		Theme:        th,
		Service:      svc,
		InvalidateFn: invalidateFn,
		Devices:      []pages.DeviceViewItem{},
		ClipItems:    []*pages.ClipItem{},
	}

	s.DeviceList.Axis = layout.Vertical
	s.ClipList.Axis = layout.Vertical
	s.SearchEditor.Submit = true
	s.IPEditor.Submit = true

	// Load pinned items from storage
	if svc != nil && svc.Storage() != nil {
		for _, text := range svc.Storage().LoadPinned() {
			s.ClipItems = append(s.ClipItems, &pages.ClipItem{
				Content:  text,
				IsPinned: true,
			})
		}
	}

	// Subscribe to event bus for real-time reactive UI updates
	if svc != nil && svc.EventBus() != nil {
		bus := svc.EventBus()

		// 1. Inbound Clip Received
		sub1 := bus.SubscribeClips(func(evt events.ClipReceivedEvent) {
			s.AddClip(evt.Content)
			if s.InvalidateFn != nil {
				s.InvalidateFn()
			}
		})
		s.subs = append(s.subs, sub1)

		// 2. Device Connected / Disconnected
		sub2 := bus.SubscribeDevices(func(evt events.DeviceEvent) {
			s.handleDeviceEvent(evt)
			if s.InvalidateFn != nil {
				s.InvalidateFn()
			}
		})
		s.subs = append(s.subs, sub2)

		// 3. Toast Notifications
		sub3 := bus.SubscribeToast(func(msg string) {
			s.TriggerToast(msg)
			if s.InvalidateFn != nil {
				s.InvalidateFn()
			}
		})
		s.subs = append(s.subs, sub3)
	}

	return s
}

func (s *AppState) TriggerToast(msg string) {
	s.ToastMsg = msg
	s.ToastStartTime = time.Now()
}

func (s *AppState) handleDeviceEvent(evt events.DeviceEvent) {
	s.DevicesMu.Lock()
	defer s.DevicesMu.Unlock()

	if evt.Type == events.DeviceConnected {
		exists := false
		for i, d := range s.Devices {
			if d.IP == evt.Device.IP {
				s.Devices[i].Name = evt.Device.Name
				exists = true
				break
			}
		}
		if !exists {
			s.Devices = append(s.Devices, pages.DeviceViewItem{
				Name: evt.Device.Name,
				IP:   evt.Device.IP,
			})
		}
	} else if evt.Type == events.DeviceDisconnected {
		var remaining []pages.DeviceViewItem
		for _, d := range s.Devices {
			if d.IP != evt.Device.IP {
				remaining = append(remaining, d)
			}
		}
		s.Devices = remaining
	}
}

func (s *AppState) AddClip(text string) {
	if text == "" {
		return
	}
	s.ClipMu.Lock()
	defer s.ClipMu.Unlock()

	// Check if already pinned
	for _, item := range s.ClipItems {
		if item.IsPinned && item.Content == text {
			return
		}
	}

	// Remove unpinned duplicate if exists
	var filtered []*pages.ClipItem
	for _, item := range s.ClipItems {
		if !item.IsPinned && item.Content == text {
			continue
		}
		filtered = append(filtered, item)
	}

	newItem := &pages.ClipItem{
		Content:  text,
		IsPinned: false,
	}

	// Insert right after pinned items
	insertIdx := 0
	for _, item := range filtered {
		if item.IsPinned {
			insertIdx++
		} else {
			break
		}
	}

	s.ClipItems = append(filtered[:insertIdx], append([]*pages.ClipItem{newItem}, filtered[insertIdx:]...)...)
}

func (s *AppState) Update(gtx layout.Context) {
	// 1. Hardware Navigation (Tabs)
	s.handleKeyEvents(gtx)

	// 2. Tab Clicks
	for i := range s.TabBtns {
		if s.TabBtns[i].Clicked(gtx) {
			s.ActiveTab = i
		}
	}

	// 3. Help Dialog
	if s.HelpBtn.Clicked(gtx) {
		s.ShowHelp = true
	}
	if s.CloseHelpBtn.Clicked(gtx) {
		s.ShowHelp = false
	}

	// 4. Manual Connect Dialog
	if s.ConnectBtn.Clicked(gtx) {
		s.ShowConnectDialog = true
	}
	if s.DialogCancelBtn.Clicked(gtx) {
		s.ShowConnectDialog = false
	}
	if s.DialogConnectBtn.Clicked(gtx) {
		ip := strings.TrimSpace(s.IPEditor.Text())
		if ip != "" && s.Service != nil {
			_ = s.Service.ConnectIP(ip)
			s.TriggerToast("Connecting to " + ip + "...")
			s.IPEditor.SetText("")
		}
		s.ShowConnectDialog = false
	}

	// 5. Clear Unpinned Clips
	if s.ClearClipsBtn.Clicked(gtx) {
		s.ClearUnpinnedClips()
	}

	// 6. Sync Phone Clip Button
	if s.SyncClipBtn.Clicked(gtx) {
		// Request clipboard content via Gio
		gtx.Execute(clipboard.ReadCmd{
			Tag: &s.SyncClipBtn,
		})
	}

	// Drain any Gio clipboard read events
	for {
		ev, ok := gtx.Event(transfer.TargetFilter{Target: &s.SyncClipBtn, Type: "text/plain"})
		if !ok {
			break
		}
		if dataEvt, ok := ev.(transfer.DataEvent); ok && dataEvt.Open != nil {
			rc := dataEvt.Open()
			if rc != nil {
				data, _ := io.ReadAll(rc)
				_ = rc.Close()
				text := string(data)
				if text != "" {
					s.AddClip(text)
					if s.Service != nil {
						s.Service.BroadcastClip(text)
					}
					s.TriggerToast("Synced & broadcast phone clip")
				}
			}
		}
	}

	// 7. Handle Individual Clip Card Actions (Copy, Pin, Delete)
	s.handleClipEvents(gtx)
}

func (s *AppState) handleClipEvents(gtx layout.Context) {
	s.ClipMu.Lock()
	defer s.ClipMu.Unlock()

	var remaining []*pages.ClipItem
	pinnedChanged := false

	for _, item := range s.ClipItems {
		// Delete action
		if item.DeleteBtn.Clicked(gtx) {
			if item.IsPinned {
				pinnedChanged = true
			}
			continue
		}

		// Pin toggle
		if item.PinBtn.Clicked(gtx) {
			item.IsPinned = !item.IsPinned
			pinnedChanged = true
			remaining = append(remaining, item)
			continue
		}

		// Expand toggle
		if item.ExpandBtn.Clicked(gtx) {
			item.IsExpanded = !item.IsExpanded
			remaining = append(remaining, item)
			continue
		}

		// Card clicked -> Copy to Android clipboard and Broadcast to peers
		if item.CardBtn.Clicked(gtx) {
			// Write to OS clipboard using Gio
			gtx.Execute(clipboard.WriteCmd{
				Type: "text/plain",
				Data: io.NopCloser(strings.NewReader(item.Content)),
			})

			// Broadcast to peers
			if s.Service != nil {
				s.Service.BroadcastClip(item.Content)
			}
			s.TriggerToast("Copied & broadcast to devices")
		}

		remaining = append(remaining, item)
	}

	s.ClipItems = remaining

	if pinnedChanged {
		s.sortClips()
		s.persistPinned()
	}
}

func (s *AppState) handleKeyEvents(gtx layout.Context) {
	if gtx.Focused(&s.SearchEditor) || gtx.Focused(&s.IPEditor) {
		for {
			ev, ok := gtx.Event(key.Filter{Name: key.NameEscape})
			if !ok {
				break
			}
			if e, ok := ev.(key.Event); ok && e.State == key.Press {
				gtx.Execute(key.FocusCmd{Tag: nil})
			}
		}
		return
	}

	for {
		ev, ok := gtx.Event(
			key.Filter{Name: key.NameLeftArrow},
			key.Filter{Name: key.NameRightArrow},
		)
		if !ok {
			break
		}
		if e, ok := ev.(key.Event); ok && e.State == key.Press {
			if e.Name == key.NameLeftArrow {
				s.ActiveTab = 0
			} else if e.Name == key.NameRightArrow {
				s.ActiveTab = 1
			}
		}
	}
}

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

func (s *AppState) ClearUnpinnedClips() {
	s.ClipMu.Lock()
	defer s.ClipMu.Unlock()
	var pinned []*pages.ClipItem
	for _, item := range s.ClipItems {
		if item.IsPinned {
			pinned = append(pinned, item)
		}
	}
	s.ClipItems = pinned
}

func (s *AppState) persistPinned() {
	if s.Service == nil || s.Service.Storage() == nil {
		return
	}
	var texts []string
	for _, item := range s.ClipItems {
		if item.IsPinned {
			texts = append(texts, item.Content)
		}
	}
	go func(saved []string) {
		_ = s.Service.Storage().SavePinned(saved)
	}(texts)
}

func (s *AppState) FilteredClips() []*pages.ClipItem {
	s.ClipMu.Lock()
	defer s.ClipMu.Unlock()

	query := strings.TrimSpace(strings.ToLower(s.SearchEditor.Text()))
	if query == "" {
		return s.ClipItems
	}

	var matched []*pages.ClipItem
	for _, item := range s.ClipItems {
		if strings.Contains(strings.ToLower(item.Content), query) || strings.Contains(strings.ToLower(utils.CleanDisplayText(item.Content)), query) {
			matched = append(matched, item)
		}
	}
	return matched
}
