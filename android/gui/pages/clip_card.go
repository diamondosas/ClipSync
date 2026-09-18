package pages

import (
	"clipsync-android/gui/themes"
	"clipsync-android/gui/utils"
	"clipsync-android/gui/widgets"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ClipItem represents an interactive mobile clipboard entry.
type ClipItem struct {
	Content    string
	IsPinned   bool
	IsExpanded bool
	CardBtn    widget.Clickable
	PinBtn     widget.Clickable
	DeleteBtn  widget.Clickable
	ExpandBtn  widget.Clickable
}

// ClipCard renders an individual clipboard history card optimized for mobile touch.
func ClipCard(gtx layout.Context, th *material.Theme, item *ClipItem) layout.Dimensions {
	bg := themes.ColorSurface
	if item.IsPinned {
		bg = themes.ColorSurfacePinned
	}

	displayText := utils.CleanDisplayText(item.Content)
	isLong := strings.Count(displayText, "\n") >= 2 || len(displayText) > 80

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return widgets.RoundedBox(gtx, 10, bg, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &item.CardBtn, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						// 1. Text display
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body2(th, displayText)
							lbl.Color = themes.ColorText
							lbl.TextSize = unit.Sp(14)
							if isLong && !item.IsExpanded {
								lbl.MaxLines = 3
								lbl.Truncator = "..."
							}
							return lbl.Layout(gtx)
						}),

						// 2. Expand/Collapse toggle for long text
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if !isLong {
								return layout.Dimensions{}
							}
							return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								label := "Show More"
								if item.IsExpanded {
									label = "Show Less"
								}
								btn := material.Button(th, &item.ExpandBtn, label)
								btn.Background = themes.ColorSurface
								btn.Color = themes.ColorBrown
								btn.TextSize = unit.Sp(12)
								btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}
								return btn.Layout(gtx)
							})
						}),

						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

						// 3. Action row: Status Tag & Pin / Delete buttons
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								// Status tag
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if item.IsPinned {
										tag := material.Caption(th, "PINNED")
										tag.Color = themes.ColorBrown
										tag.TextSize = unit.Sp(11)
										return tag.Layout(gtx)
									}
									return layout.Dimensions{}
								}),

								layout.Flexed(1, layout.Spacer{}.Layout),

								// Action buttons
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										// Pin / Unpin button
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											pinLabel := "📌 Pin"
											pinBg := themes.ColorBg
											pinFg := themes.ColorTextMuted
											if item.IsPinned {
												pinLabel = "📌 Pinned"
												pinBg = themes.ColorBrown
												pinFg = themes.ColorBg
											}
											btn := material.Button(th, &item.PinBtn, pinLabel)
											btn.Background = pinBg
											btn.Color = pinFg
											btn.TextSize = unit.Sp(12)
											btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(10), Right: unit.Dp(10)}
											return btn.Layout(gtx)
										}),

										layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),

										// Delete button
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btn := material.Button(th, &item.DeleteBtn, "✕")
											btn.Background = themes.ColorBg
											btn.Color = themes.ColorRed
											btn.TextSize = unit.Sp(12)
											btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(10), Right: unit.Dp(10)}
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
