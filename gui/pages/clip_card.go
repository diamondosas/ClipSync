// ai-generated
package pages

import (
	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ClipItem represents an interactive clipboard entry.
type ClipItem struct {
	Content   string
	IsPinned  bool
	CardBtn   widget.Clickable
	PinBtn    widget.Clickable
	DeleteBtn widget.Clickable
}

// ClipCard renders an individual clipboard history card with copy, pin, and delete actions.
func ClipCard(gtx layout.Context, th *material.Theme, item *ClipItem) layout.Dimensions {
	bg := themes.ColorSurface
	if item.IsPinned {
		bg = themes.ColorSurfacePinned
	}

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 8, bg, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// 1. Clickable text area for copying
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Clickable(gtx, &item.CardBtn, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body2(th, item.Content)
							lbl.Color = themes.ColorText
							lbl.MaxLines = 3
							lbl.Truncator = "..."
							return lbl.Layout(gtx)
						})
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
					// 2. Action row: Status label + [Pin] + [X]
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
							// Status hint
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if item.IsPinned {
									tag := material.Caption(th, "PINNED")
									tag.Color = themes.ColorCyan
									return tag.Layout(gtx)
								}
								hint := material.Caption(th, "Click to copy")
								hint.Color = themes.ColorTextMuted
								return hint.Layout(gtx)
							}),
							// Buttons on right
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									// Pin button
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										pinLabel := "Pin"
										pinBg := themes.ColorSurface
										pinFg := themes.ColorTextMuted
										if item.IsPinned {
											pinLabel = "Unpin"
											pinBg = themes.ColorCyan
											pinFg = themes.ColorBg
										}
										btn := material.Button(th, &item.PinBtn, pinLabel)
										btn.Background = pinBg
										btn.Color = pinFg
										btn.TextSize = unit.Sp(11)
										btn.Inset = layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3), Left: unit.Dp(6), Right: unit.Dp(6)}
										return btn.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
									// Delete button
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, &item.DeleteBtn, "X")
										btn.Background = themes.ColorSurface
										btn.Color = themes.ColorRed
										btn.TextSize = unit.Sp(11)
										btn.Inset = layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3), Left: unit.Dp(6), Right: unit.Dp(6)}
										return btn.Layout(gtx)
									}),
								)
							}),
						)
					}),
				)
			})
		})
	})
}
