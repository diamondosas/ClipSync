package main

import (
	"log"
	"strconv"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		// No os.Exit here: it would kill the whole Java app when the
		// Gio screen closes. Looping lets the screen be opened again.
		for {
			w := new(app.Window)
			if err := run(w); err != nil {
				log.Println("gio window closed:", err)
			}
		}
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	var (
		ops   op.Ops
		btn   widget.Clickable
		count int
	)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if btn.Clicked(gtx) {
				count++
			}
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(material.H4(th, "Hello from Gio").Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
					layout.Rigid(material.Button(th, &btn, "Clicked "+strconv.Itoa(count)).Layout),
				)
			})
			e.Frame(gtx.Ops)
		}
	}
}

