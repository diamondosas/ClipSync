package components

import (
	"image/color"

	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)
var helpMsg string =  "Start the Application on another device & Make sure devices are on the same network."
// HelpDialog overlays a help message over the current layout if show is true.
func HelpDialog(gtx layout.Context, th *material.Theme, closeBtn *widget.Clickable, show bool, underlying layout.Widget) layout.Dimensions {
	// Always lay out the underlying content first
	dims := underlying(gtx)

	if !show {
		return dims
	}

	// Semi-transparent overlay color
	scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: 180}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		// 1. Background content (already rendered, we just need to stack it so it's behind)
		// Actually, standard Gio way is to render the underlying content inside the stack,
		// but since we called it above to get dims, we can just render it again or assume
		// the overlay draws on top of whatever was just drawn. 
		// By drawing our scrim with `Expanded` it will cover the previous drawing operations.
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return widgets.ColorBox(gtx, scrimColor, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		// 2. Dialog box centered
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return widgets.RoundedBox(gtx, 10, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							// Dialog Title
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, "Help & Info")
								title.Color = themes.ColorCyan
								return title.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
							// Dialog Body text
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								body := material.Body1(th, helpMsg)
								body.Color = themes.ColorText
								body.Alignment = text.Middle
								return body.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
							// OK Button to close
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, closeBtn, "OK")
								btn.Background = themes.ColorCyan
								btn.Color = themes.ColorBg
								return btn.Layout(gtx)
							}),
						)
					})
				})
			})
		}),
	)
}
