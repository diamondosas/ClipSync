// ai-generated
package pages

import (
	"fmt"

	"clipsync/gui/themes"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ClipboardPage lays out the clipboard history list with clear and interactive actions.
func ClipboardPage(gtx layout.Context, th *material.Theme, list *widget.List, clips []*ClipItem, clearBtn *widget.Clickable) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Sub-header bar: Item Count & Clear All Button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						countStr := fmt.Sprintf("%d clips", len(clips))
						if len(clips) == 1 {
							countStr = "1 clip"
						}
						lbl := material.Caption(th, countStr)
						lbl.Color = themes.ColorTextMuted
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if len(clips) == 0 {
							return layout.Dimensions{}
						}
						btn := material.Button(th, clearBtn, "Clear All")
						btn.Background = themes.ColorSurface
						btn.Color = themes.ColorTextMuted
						btn.TextSize = unit.Sp(11)
						btn.Inset = layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3), Left: unit.Dp(8), Right: unit.Dp(8)}
						return btn.Layout(gtx)
					}),
				)
			})
		}),
		// Scrollable List or Empty State
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if len(clips) == 0 {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(th, "Clipboard is empty.\nCopy text on any device to see it here.")
						lbl.Color = themes.ColorTextMuted
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					})
				})
			}
			return list.Layout(gtx, len(clips), func(gtx layout.Context, index int) layout.Dimensions {
				return ClipCard(gtx, th, clips[index])
			})
		}),
	)
}
