package gui

import (
	"clipsync-android/gui/components"
	"clipsync-android/gui/pages"
	"clipsync-android/gui/themes"
	"clipsync-android/internal/service"
	"log"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var (
	MainWindow *app.Window
)

// StartGUI initializes and runs the Gio-based mobile user interface.
func StartGUI(svc *service.ClipSyncService) {
	go func() {
		w := new(app.Window)
		MainWindow = w

		if err := run(w, svc); err != nil {
			log.Printf("[GUI] Window event loop error: %v", err)
		}
	}()
	app.Main()
}

func run(w *app.Window, svc *service.ClipSyncService) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette.Fg = themes.ColorText
	th.Palette.Bg = themes.ColorBg
	th.Palette.ContrastBg = themes.ColorBrown
	th.Palette.ContrastFg = themes.ColorBg

	state := NewAppState(th, svc, func() {
		if MainWindow != nil {
			MainWindow.Invalidate()
		}
	})

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Process events (clicks, gestures, inputs)
			state.Update(gtx)

			// Layout the mobile interface
			layoutMain(gtx, state)

			e.Frame(gtx.Ops)
		}
	}
}

func layoutMain(gtx layout.Context, s *AppState) layout.Dimensions {
	// Fill full mobile canvas background directly
	paint.FillShape(gtx.Ops, themes.ColorBg, clip.Rect{Max: gtx.Constraints.Max}.Op())

	mainContent := func(gtx layout.Context) layout.Dimensions {
		event.Op(gtx.Ops, s)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		s.DevicesMu.Lock()
		connectedCount := len(s.Devices)
		s.DevicesMu.Unlock()

		dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// 1. Mobile Top AppBar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := "ClipSync"
				if s.ActiveTab == 1 {
					title = "Clipboard"
				}
				return components.AppBar(gtx, s.Theme, title, connectedCount, &s.HelpBtn, &s.ConnectBtn)
			}),

			// 2. Dynamic Page Body
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if s.ActiveTab == 0 {
					s.DevicesMu.Lock()
					devList := make([]pages.DeviceViewItem, len(s.Devices))
					copy(devList, s.Devices)
					s.DevicesMu.Unlock()

					return pages.DevicesPage(gtx, s.Theme, &s.DeviceList, devList, &s.ConnectBtn)
				}

				s.ClipMu.Lock()
				totalClips := len(s.ClipItems)
				s.ClipMu.Unlock()

				return pages.ClipboardPage(
					gtx,
					s.Theme,
					&s.ClipList,
					s.FilteredClips(),
					&s.ClearClipsBtn,
					&s.SyncClipBtn,
					&s.SearchEditor,
					&s.SearchClearBtn,
					totalClips,
				)
			}),

			// 3. Mobile Bottom Navigation
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return components.BottomNav(gtx, s.Theme, s.ActiveTab, &s.TabBtns)
			}),
		)

		// Floating Toast notification banner
		if s.ToastMsg != "" && !s.ToastStartTime.IsZero() {
			layout.Stack{Alignment: layout.S}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Bottom: unit.Dp(64)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return components.Toast(gtx, s.Theme, s.ToastMsg, s.ToastStartTime)
					})
				}),
			)
		}

		return dims
	}

	// Overlay Connect Dialog if active
	afterConnect := func(gtx layout.Context) layout.Dimensions {
		return components.ConnectDialog(gtx, s.Theme, &s.IPEditor, &s.DialogConnectBtn, &s.DialogCancelBtn, s.ShowConnectDialog, mainContent)
	}

	// Overlay Help Dialog if active
	return components.HelpDialog(gtx, s.Theme, &s.CloseHelpBtn, s.ShowHelp, afterConnect)
}
