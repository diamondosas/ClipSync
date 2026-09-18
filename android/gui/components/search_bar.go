package components

import (
	"clipsync-android/gui/themes"
	"clipsync-android/gui/widgets"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// SearchBar renders a clean text search input for mobile with an optional clear button.
func SearchBar(gtx layout.Context, th *material.Theme, editor *widget.Editor, clearBtn *widget.Clickable) layout.Dimensions {
	editor.Submit = true

	// Clean accidental newlines or carriage returns
	if txt := editor.Text(); strings.ContainsAny(txt, "\n\r\t") {
		cleaned := strings.ReplaceAll(txt, "\t", " ")
		cleaned = strings.ReplaceAll(strings.ReplaceAll(cleaned, "\r", ""), "\n", "")
		editor.SetText(cleaned)
	}

	if clearBtn.Clicked(gtx) {
		editor.SetText("")
	}

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 8, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
			return widgets.TouchTarget(gtx, unit.Dp(44), func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							ed := material.Editor(th, editor, "Search clips...")
							ed.Color = themes.ColorText
							ed.HintColor = themes.ColorTextMuted
							ed.TextSize = unit.Sp(14)
							return ed.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if len(editor.Text()) == 0 {
								return layout.Dimensions{}
							}
							btn := material.Button(th, clearBtn, "✕")
							btn.Background = themes.ColorSurface
							btn.Color = themes.ColorTextMuted
							btn.TextSize = unit.Sp(11)
							btn.Inset = layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(6), Right: unit.Dp(6)}
							return btn.Layout(gtx)
						}),
					)
				})
			})
		})
	})
}
