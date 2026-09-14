// ai-generated
package components

import (
	"image"
	"image/color"

	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// TitleBar renders a sleek, cross-platform custom window header with draggable
// moving area, title, help button, minimize button, and close button.
// It completely replaces default OS/Gio client-side decorations (no purple bar, no maximize button).
func TitleBar(
	gtx layout.Context,
	th *material.Theme,
	helpBtn *widget.Clickable,
	minBtn *widget.Clickable,
	closeBtn *widget.Clickable,
	titleText string,
) layout.Dimensions {
	// Top title bar container with bottom margin separating it from page elements
	return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.ColorBox(gtx, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(12),
				Right:  unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					// 1. Left side: App title (Draggable)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							// Small decorative accent dot
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return widgets.RoundedBox(gtx, 4, themes.ColorBrown, func(gtx layout.Context) layout.Dimensions {
									return layout.Dimensions{Size: image.Pt(gtx.Dp(unit.Dp(8)), gtx.Dp(unit.Dp(8)))}
								})
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							// Title text
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, titleText)
								title.Color = themes.ColorBrown
								title.TextSize = unit.Sp(16)
								return title.Layout(gtx)
							}),
						)

						// Allow dragging from title area
						defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()
						system.ActionInputOp(system.ActionMove).Add(gtx.Ops)
						return dims
					}),

					// 2. Middle area: Draggable spacer (click & drag anywhere to move window)
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
					}),

					// 3. Right side: Action controls [ ? ] [ – ] [ ✕ ]
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							// Help button "?"
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return windowButton(gtx, th, helpBtn, "?", themes.ColorTextMuted, unit.Sp(12))
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),

							// Minimize button "–"
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return windowButton(gtx, th, minBtn, "–", themes.ColorTextMuted, unit.Sp(12))
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),

							// Close button "✕" (with subtle red hover contrast)
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								fg := themes.ColorTextMuted
								if closeBtn.Hovered() {
									fg = themes.ColorRed
								}
								return windowButton(gtx, th, closeBtn, "✕", fg, unit.Sp(11))
							}),
						)
					}),
				)
			})
		})
	})
}

// windowButton renders an individual sleek, rounded window titlebar button.
func windowButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, symbol string, fg color.NRGBA, textSize unit.Sp) layout.Dimensions {
	return material.Clickable(gtx, btn, func(gtx layout.Context) layout.Dimensions {
		bg := themes.ColorBg
		if btn.Hovered() {
			bg = themes.ColorSurfacePinned
		}
		return widgets.RoundedBox(gtx, 4, bg, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(4),
				Bottom: unit.Dp(4),
				Left:   unit.Dp(7),
				Right:  unit.Dp(7),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(th, symbol)
				lbl.Color = fg
				lbl.TextSize = textSize
				lbl.Alignment = text.Middle
				return lbl.Layout(gtx)
			})
		})
	})
}
