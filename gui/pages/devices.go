package pages

import (
	"fmt"

	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Device struct {
	Name string
	IP   string
}

// DevicesPage lays out the connection info and devices list.
func DevicesPage(gtx layout.Context, th *material.Theme, list *widget.List, devices []Device) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Sub-header text
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(12), Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, "looking for devices...")
				lbl.Color = themes.ColorTextMuted
				return lbl.Layout(gtx)
			})
		}),
		// Scrollable List of Devices
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return list.Layout(gtx, len(devices), func(gtx layout.Context, index int) layout.Dimensions {
				return deviceCard(gtx, th, devices[index])
			})
		}),
	)
}

// deviceCard renders an individual device card in the list.
func deviceCard(gtx layout.Context, th *material.Theme, dev Device) layout.Dimensions {
	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 8, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						name := material.Body1(th, dev.Name)
						name.Color = themes.ColorText
						return name.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						ip := material.Caption(th, fmt.Sprintf("IP: %s", dev.IP))
						ip.Color = themes.ColorTextMuted
						return ip.Layout(gtx)
					}),
				)
			})
		})
	})
}
