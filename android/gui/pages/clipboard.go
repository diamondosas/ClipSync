package pages

import (
	"clipsync-android/gui/components"
	"clipsync-android/gui/themes"
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ClipboardPage lays out the search bar, action subheader, and clip list.
func ClipboardPage(
	gtx layout.Context,
	th *material.Theme,
	list *widget.List,
	clips []*ClipItem,
	clearBtn *widget.Clickable,
	syncBtn *widget.Clickable,
	searchEditor *widget.Editor,
	searchClearBtn *widget.Clickable,
	totalClips int,
) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 1. Search Bar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.SearchBar(gtx, th, searchEditor, searchClearBtn)
		}),

		// 2. Subheader Bar (Item Count & Action Buttons)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						countStr := fmt.Sprintf("%d clips", totalClips)
						if totalClips == 1 {
							countStr = "1 clip"
						}
						if len(searchEditor.Text()) > 0 {
							countStr = fmt.Sprintf("%d of %d clips", len(clips), totalClips)
						}
						lbl := material.Caption(th, countStr)
						lbl.Color = themes.ColorTextMuted
						lbl.TextSize = unit.Sp(12)
						return lbl.Layout(gtx)
					}),

					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							// Quick "Sync Clip" button
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if syncBtn == nil {
									return layout.Dimensions{}
								}
								btn := material.Button(th, syncBtn, "Sync Phone Clip")
								btn.Background = themes.ColorSurface
								btn.Color = themes.ColorBrown
								btn.TextSize = unit.Sp(11)
								btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}
								return btn.Layout(gtx)
							}),

							layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),

							// Clear All Unpinned button
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if totalClips == 0 {
									return layout.Dimensions{}
								}
								btn := material.Button(th, clearBtn, "Clear All")
								btn.Background = themes.ColorSurface
								btn.Color = themes.ColorTextMuted
								btn.TextSize = unit.Sp(11)
								btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}
								return btn.Layout(gtx)
							}),
						)
					}),
				)
			})
		}),

		// 3. Scrollable List or Empty State
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if totalClips == 0 {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(th, "Clipboard history is empty.\nCopy text or receive clips from connected devices.")
						lbl.Color = themes.ColorTextMuted
						lbl.Alignment = text.Middle
						lbl.TextSize = unit.Sp(14)
						return lbl.Layout(gtx)
					})
				})
			}

			if len(clips) == 0 {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(th, fmt.Sprintf("No clips matching %q", searchEditor.Text()))
						lbl.Color = themes.ColorTextMuted
						lbl.Alignment = text.Middle
						lbl.TextSize = unit.Sp(14)
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
