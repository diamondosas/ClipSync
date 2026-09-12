// ai-generated
package components

import (
	"strings"

	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// SearchBar renders a clean text search input with a clear button following learn.go pattern.
func SearchBar(gtx layout.Context, th *material.Theme, editor *widget.Editor, clearBtn *widget.Clickable) layout.Dimensions {
	editor.Submit = true

	// Clean any accidental newlines so input stays on one straight line
	if txt := editor.Text(); strings.Contains(txt, "\n") || strings.Contains(txt, "\r") {
		cleaned := strings.ReplaceAll(strings.ReplaceAll(txt, "\r", ""), "\n", "")
		editor.SetText(cleaned)
	}

	// If clear button is clicked, empty the search text
	if clearBtn.Clicked(gtx) {
		editor.SetText("")
	}

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 6, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					// Search input editor (learn.go pattern)
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						ed := material.Editor(th, editor, "Search clips...")
						ed.Color = themes.ColorText
						ed.HintColor = themes.ColorTextMuted
						ed.TextSize = unit.Sp(13)
						return ed.Layout(gtx)
					}),
					// Clear button (only shown when query is not empty)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if len(editor.Text()) == 0 {
							return layout.Dimensions{}
						}
						btn := material.Button(th, clearBtn, "X")
						btn.Background = themes.ColorSurface
						btn.Color = themes.ColorTextMuted
						btn.TextSize = unit.Sp(10)
						btn.Inset = layout.Inset{Top: unit.Dp(1), Bottom: unit.Dp(1), Left: unit.Dp(6), Right: unit.Dp(6)}
						return btn.Layout(gtx)
					}),
				)
			})
		})
	})
}
