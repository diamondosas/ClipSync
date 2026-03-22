package gui

import (
	"image/color"
	"log"
	"os"

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
			app.Title("🚀 ClipSync"),
			app.Size(unit.Dp(400), unit.Dp(300)),
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
	
	// Use the "Electric Cyan" color from the palette
	electricCyan := color.NRGBA{R: 0x00, G: 0xC2, B: 0xFF, A: 0xFF}
	th.Palette.ContrastBg = electricCyan

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Simple vertical layout with a centered label
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H4(th, "🚀 ClipSync")
						title.Color = electricCyan
						title.Alignment = text.Middle
						return title.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(th, "Status: Running")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						caption := material.Caption(th, "Syncing clipboard between your devices.")
						caption.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
						caption.Alignment = text.Middle
						return caption.Layout(gtx)
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}
