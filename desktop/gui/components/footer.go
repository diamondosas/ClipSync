package components

import (
	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// FooterTabs lays out the bottom navigation bar with tabs.
func FooterTabs(gtx layout.Context, th *material.Theme, activeTab int, tabBtns *[2]widget.Clickable) layout.Dimensions {
	labels := []string{"Devices", "Clipboard"}

	return widgets.ColorBox(gtx, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return tabButton(gtx, th, &tabBtns[0], 0, activeTab, labels[0])
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return tabButton(gtx, th, &tabBtns[1], 1, activeTab, labels[1])
			}),
		)
	})
}

// tabButton renders an individual tab button.
func tabButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, index, activeIndex int, label string) layout.Dimensions {
	active := index == activeIndex
	bg := themes.ColorSurface
	fg := themes.ColorTextMuted

	if active {
		bg = themes.ColorCyan
		fg = themes.ColorBg // Dark text on bright cyan background for contrast
	}

	return material.Clickable(gtx, btn, func(gtx layout.Context) layout.Dimensions {
		return widgets.ColorBox(gtx, bg, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, label)
				lbl.Color = fg
				lbl.Alignment = text.Middle
				return lbl.Layout(gtx)
			})
		})
	})
}
