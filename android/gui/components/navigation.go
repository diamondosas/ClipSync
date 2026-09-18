package components

import (
	"clipsync-android/gui/themes"
	"clipsync-android/gui/widgets"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// BottomNav renders the mobile bottom navigation bar with Devices and Clipboard tabs.
func BottomNav(gtx layout.Context, th *material.Theme, activeTab int, tabBtns *[2]widget.Clickable) layout.Dimensions {
	labels := []string{"Devices", "Clipboard"}

	return widgets.ColorBox(gtx, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return navButton(gtx, th, &tabBtns[0], 0, activeTab, labels[0])
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return navButton(gtx, th, &tabBtns[1], 1, activeTab, labels[1])
			}),
		)
	})
}

func navButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, index, activeIndex int, label string) layout.Dimensions {
	active := (index == activeIndex)
	bg := themes.ColorSurface
	fg := themes.ColorTextMuted

	if active {
		bg = themes.ColorBrown
		fg = themes.ColorBg
	}

	return material.Clickable(gtx, btn, func(gtx layout.Context) layout.Dimensions {
		return widgets.ColorBox(gtx, bg, func(gtx layout.Context) layout.Dimensions {
			return widgets.TouchTarget(gtx, unit.Dp(48), func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(14)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th, label)
					lbl.Color = fg
					lbl.Alignment = text.Middle
					lbl.TextSize = unit.Sp(14)
					return lbl.Layout(gtx)
				})
			})
		})
	})
}
