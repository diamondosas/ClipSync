package components

import (
	"clipsync-android/gui/themes"
	"clipsync-android/gui/widgets"
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var helpMsg = "ClipSync synchronizes clipboard data across your local Wi-Fi devices without any cloud servers.\n\n" +
	"• Tap any clip card to copy and broadcast it to all connected peers.\n" +
	"• Use '+ IP' to manually connect if Wi-Fi blocks automatic discovery.\n" +
	"• Tap the Pin icon to preserve important snippets."

// HelpDialog overlays a modal with information when show is true.
func HelpDialog(gtx layout.Context, th *material.Theme, closeBtn *widget.Clickable, show bool, underlying layout.Widget) layout.Dimensions {
	dims := underlying(gtx)
	if !show {
		return dims
	}

	scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: 190}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return widgets.ColorBox(gtx, scrimColor, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return widgets.RoundedBox(gtx, 12, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, "About ClipSync")
								title.Color = themes.ColorBrown
								title.TextSize = unit.Sp(18)
								return title.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								body := material.Body2(th, helpMsg)
								body.Color = themes.ColorText
								body.Alignment = text.Start
								return body.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(18)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, closeBtn, "Got it")
								btn.Background = themes.ColorBrown
								btn.Color = themes.ColorBg
								btn.TextSize = unit.Sp(14)
								btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(24), Right: unit.Dp(24)}
								return btn.Layout(gtx)
							}),
						)
					})
				})
			})
		}),
	)
}

// ConnectDialog overlays a modal to enter an IP address manually.
func ConnectDialog(
	gtx layout.Context,
	th *material.Theme,
	ipEditor *widget.Editor,
	connectBtn *widget.Clickable,
	cancelBtn *widget.Clickable,
	show bool,
	underlying layout.Widget,
) layout.Dimensions {
	dims := underlying(gtx)
	if !show {
		return dims
	}

	scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: 190}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return widgets.ColorBox(gtx, scrimColor, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return widgets.RoundedBox(gtx, 12, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, "Connect Device IP")
								title.Color = themes.ColorBrown
								title.TextSize = unit.Sp(18)
								return title.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return widgets.RoundedBox(gtx, 6, themes.ColorBg, func(gtx layout.Context) layout.Dimensions {
									return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										ed := material.Editor(th, ipEditor, "e.g. 192.168.1.50")
										ed.Color = themes.ColorText
										ed.HintColor = themes.ColorTextMuted
										ed.TextSize = unit.Sp(14)
										return ed.Layout(gtx)
									})
								})
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(18)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, cancelBtn, "Cancel")
										btn.Background = themes.ColorSurface
										btn.Color = themes.ColorTextMuted
										btn.TextSize = unit.Sp(13)
										btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}
										return btn.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, connectBtn, "Connect")
										btn.Background = themes.ColorBrown
										btn.Color = themes.ColorBg
										btn.TextSize = unit.Sp(13)
										btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}
										return btn.Layout(gtx)
									}),
								)
							}),
						)
					})
				})
			})
		}),
	)
}
