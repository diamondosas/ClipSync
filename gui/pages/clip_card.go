// ai-generated
package pages

import (
	"strings"

	"clipsync/gui/assets"
	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ClipItem represents an interactive clipboard entry.
type ClipItem struct {
	Content    string
	IsPinned   bool
	IsExpanded bool
	CardBtn    widget.Clickable
	PinBtn     widget.Clickable
	DeleteBtn  widget.Clickable
	ExpandBtn  widget.Clickable
}

// ClipCard renders an individual clipboard history card with copy, pin, delete, and smart preview.
func ClipCard(gtx layout.Context, th *material.Theme, item *ClipItem) layout.Dimensions {
	bg := themes.ColorSurface
	if item.IsPinned {
		bg = themes.ColorSurfacePinned
	}

	isLong := strings.Count(item.Content, "\n") >= 2 || len(item.Content) > 80

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return widgets.RoundedBox(gtx, 8, bg, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &item.CardBtn, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						// 1. Text display (entire card is clickable)
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body2(th, item.Content)
							lbl.Color = themes.ColorText
							if isLong && !item.IsExpanded {
								lbl.MaxLines = 2
								lbl.Truncator = "..."
							}
							return lbl.Layout(gtx)
						}),
						// 2. Expand/Collapse toggle for long text clips
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if !isLong {
								return layout.Dimensions{}
							}
							return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								label := "More"
								if item.IsExpanded {
									label = "Less"
								}
								btn := material.Button(th, &item.ExpandBtn, label)
								btn.Background = themes.ColorSurface
								btn.Color = themes.ColorCyan
								btn.TextSize = unit.Sp(10)
								btn.Inset = layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(6), Right: unit.Dp(6)}
								return btn.Layout(gtx)
							})
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
						// 3. Action row: Status label (PINNED if pinned) + [Pin] + [X]
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								// Status hint: PINNED tag if pinned, otherwise clean empty
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if item.IsPinned {
										tag := material.Caption(th, "PINNED")
										tag.Color = themes.ColorCyan
										return tag.Layout(gtx)
									}
									return layout.Dimensions{}
								}),
								// Flexible spacer pushing buttons to the far right
								layout.Flexed(1, layout.Spacer{}.Layout),
								// Buttons on right
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										// Pin / Unpin button with pushpin image icon
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											pinBg := themes.ColorSurface
											if item.IsPinned {
												pinBg = themes.ColorCyan
											}

											return material.Clickable(gtx, &item.PinBtn, func(gtx layout.Context) layout.Dimensions {
												return widgets.RoundedBox(gtx, 4, pinBg, func(gtx layout.Context) layout.Dimensions {
													return layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3), Left: unit.Dp(6), Right: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
														img := widget.Image{
															Src:      assets.PinImageOp,
															Scale:    0.6,
															Position: layout.Center,
														}
														return img.Layout(gtx)
													})
												})
											})
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
	})
}
