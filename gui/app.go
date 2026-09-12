// ai-generated
package gui

import (
	_ "embed"
	_ "image/jpeg"
	"log"
	"sync"

	// "os"

	"clipsync/gui/components"
	"clipsync/gui/pages"
	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)
var(
	Window *app.Window
	windowMu sync.Mutex
	isWindowOpen bool
)
// StartGUI initializes and runs the Gio-based user interface.
// This function will block until the application is closed.
func StartGUI() {
	ShowWindow()
	app.Main()
}


func ShowWindow(){
	windowMu.Lock()
	defer windowMu.Unlock()

	if isWindowOpen {
		if Window != nil{
			Window.Perform(system.ActionRaise)
		}
		return
	}
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("ClipSync"),
			app.Size(unit.Dp(300), unit.Dp(480)),
			app.MinSize(unit.Dp(260), unit.Dp(380)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}

		windowMu.Lock()
		isWindowOpen = false
		Window = nil
		windowMu.Unlock()
		log.Println("[GUI] Window closed. ClipSync continues running in system tray.")

	}()
}
func run(w *app.Window) error {
	// Initialize a material theme with default fonts
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette.Fg = themes.ColorText
	th.Palette.Bg = themes.ColorBg
	th.Palette.ContrastBg = themes.ColorBrown
	th.Palette.ContrastFg = themes.ColorBg
	Window = w
	// Create application state
	state := NewAppState(th)

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Process events (clicks, inputs)
			state.Update(gtx)

			// Layout the entire window
			layoutMain(gtx, state)

			e.Frame(gtx.Ops)
		}
	}
}

// layoutMain is the root layout builder for the application.
func layoutMain(gtx layout.Context, s *AppState) layout.Dimensions {
	// Root background color
	return widgets.ColorBox(gtx, themes.ColorBg, func(gtx layout.Context) layout.Dimensions {

		// The main content of the application (Header, Body, Footer + Toast)
		mainContent := func(gtx layout.Context) layout.Dimensions {
			// Register for global keyboard events
			event.Op(gtx.Ops, s)

			// Enforce full width so SpaceBetween and right-aligned buttons spread across the window
			gtx.Constraints.Min.X = gtx.Constraints.Max.X

			dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// Header (Persistent)
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := "ClipSync"
					if s.ActiveTab == 1 {
						title = "Clipboard"
					}
					return components.Header(gtx, s.Theme, &s.HelpBtn, title)
				}),

				// Body (Dynamic Page Content)
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if s.ActiveTab == 0 {
						return pages.DevicesPage(gtx, s.Theme, &s.DeviceList, s.Devices)
					}
					return pages.ClipboardPage(
						gtx,
						s.Theme,
						&s.ClipList,
						s.FilteredClips(),
						&s.ClearClipsBtn,
						&s.SearchEditor,
						&s.SearchClearBtn,
						len(s.ClipItems),
					)
				}),

				// Footer (Navigation Tabs)
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return components.FooterTabs(gtx, s.Theme, s.ActiveTab, &s.TabBtns)
				}),
			)

			// Floating Toast Notification (rendered on top above footer)
			if s.ToastMsg != "" && !s.ToastStartTime.IsZero() {
				layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Bottom: unit.Dp(56)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return components.Toast(gtx, s.Theme, s.ToastMsg, s.ToastStartTime)
						})
					}),
				)
			}

			return dims
		}

		// Overlay the Help Dialog if `s.ShowHelp` is true
		return components.HelpDialog(gtx, s.Theme, &s.CloseHelpBtn, s.ShowHelp, mainContent)
	})
}
