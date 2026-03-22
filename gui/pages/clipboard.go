package pages

import (
	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ClipboardPage lays out the clipboard history list.
func ClipboardPage(gtx layout.Context, th *material.Theme, list *widget.List, history []string) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Scrollable List of Clipboard items
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return list.Layout(gtx, len(history), func(gtx layout.Context, index int) layout.Dimensions {
				return clipboardCard(gtx, th, history[index])
			})
		}),
	)
}

// clipboardCard renders an individual clipboard content item.
func clipboardCard(gtx layout.Context, th *material.Theme, content string) layout.Dimensions {
	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 8, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, content)
				lbl.Color = themes.ColorText
				return lbl.Layout(gtx)
			})
		})
	})
}
