package components

import (
	"clipsync/gui/themes"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Header lays out the top bar with Title and Help button.
func Header(gtx layout.Context, th *material.Theme, helpBtn *widget.Clickable, titleText string) layout.Dimensions {
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			// Left side: Title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(th, titleText)
				title.Color = themes.ColorCyan
				return title.Layout(gtx)
			}),
			// Right side: Help Button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, helpBtn, "?")
				btn.Background = themes.ColorSurface
				btn.Color = themes.ColorText
				btn.Inset = layout.UniformInset(unit.Dp(8))
				return btn.Layout(gtx)
			}),
		)
	})
}
