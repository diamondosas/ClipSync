package gui

import (
	_ "embed"
	_ "image/jpeg"
	"log"
	"os"

	"clipsync/gui/components"
	"clipsync/gui/pages"
	"clipsync/gui/themes"
	"clipsync/gui/widgets" 

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// StartGUI initializes and runs the Gio-based user interface.
// This function will block until the application is closed.
func StartGUI() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("ClipSync"),
			app.MaxSize(unit.Dp(250), unit.Dp(400)),
			app.MinSize(unit.Dp(250), unit.Dp(400)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	// Initialize a material theme with default fonts
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

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

		// The main content of the application (Header, Body, Footer)
		mainContent := func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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
					return pages.ClipboardPage(gtx, s.Theme, &s.ClipList, s.History)
				}),

				// Footer (Navigation Tabs)
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return components.FooterTabs(gtx, s.Theme, s.ActiveTab, &s.TabBtns)
				}),
			)
		}

		// Overlay the Help Dialog if `s.ShowHelp` is true
		return components.HelpDialog(gtx, s.Theme, &s.CloseHelpBtn, s.ShowHelp, mainContent)
	})
}
