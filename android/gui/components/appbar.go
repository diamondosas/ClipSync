package components

import (
	"clipsync-android/gui/themes"
	"clipsync-android/gui/widgets"
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// AppBar renders the mobile top app bar with title, connection status indicator, and action buttons.
func AppBar(
	gtx layout.Context,
	th *material.Theme,
	title string,
	connectedCount int,
	helpBtn *widget.Clickable,
	connectBtn *widget.Clickable,
) layout.Dimensions {
	return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.ColorBox(gtx, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(12),
				Bottom: unit.Dp(12),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					// Left side: Title and connected status
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							// Status Dot
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								dotColor := themes.ColorTextMuted
								if connectedCount > 0 {
									dotColor = themes.ColorGreen
								}
								return widgets.RoundedBox(gtx, 5, dotColor, func(gtx layout.Context) layout.Dimensions {
									return layout.Dimensions{Size: image.Pt(gtx.Dp(unit.Dp(10)), gtx.Dp(unit.Dp(10)))}
								})
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							// Title Text
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								titleLabel := material.H6(th, title)
								titleLabel.Color = themes.ColorBrown
								titleLabel.TextSize = unit.Sp(18)
								return titleLabel.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							// Connected count tag
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if connectedCount == 0 {
									return layout.Dimensions{}
								}
								tagText := fmt.Sprintf("(%d online)", connectedCount)
								tag := material.Caption(th, tagText)
								tag.Color = themes.ColorTextMuted
								return tag.Layout(gtx)
							}),
						)
					}),

					// Right side: Action Buttons (Connect IP & Help)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							// Connect manual IP button "+"
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if connectBtn == nil {
									return layout.Dimensions{}
								}
								btn := material.Button(th, connectBtn, "+ IP")
								btn.Background = themes.ColorBg
								btn.Color = themes.ColorBrown
								btn.TextSize = unit.Sp(12)
								btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}
								return btn.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							// Help button "?"
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Clickable(gtx, helpBtn, func(gtx layout.Context) layout.Dimensions {
									return widgets.RoundedBox(gtx, 4, themes.ColorBg, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											lbl := material.Caption(th, "?")
											lbl.Color = themes.ColorTextMuted
											lbl.TextSize = unit.Sp(13)
											lbl.Alignment = text.Middle
											return lbl.Layout(gtx)
										})
									})
								})
							}),
						)
					}),
				)
			})
		})
	})
}
