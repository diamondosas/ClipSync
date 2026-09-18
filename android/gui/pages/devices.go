package pages

import (
	"clipsync-android/gui/themes"
	"clipsync-android/gui/widgets"
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type DeviceViewItem struct {
	Name string
	IP   string
}

// DevicesPage lays out the discovered devices list and search status.
func DevicesPage(
	gtx layout.Context,
	th *material.Theme,
	list *widget.List,
	devices []DeviceViewItem,
	connectBtn *widget.Clickable,
) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Sub-header with animated searching dots
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			now := gtx.Now
			if now.IsZero() {
				now = time.Now()
			}
			nowMillis := now.UnixMilli()
			interval := int64(450)
			step := (nowMillis / interval) % 4
			dots := ""
			switch step {
			case 1:
				dots = "."
			case 2:
				dots = ".."
			case 3:
				dots = "..."
			}

			remaining := interval - (nowMillis % interval)
			if remaining <= 0 {
				remaining = interval
			}
			gtx.Execute(op.InvalidateCmd{At: now.Add(time.Duration(remaining) * time.Millisecond)})

			return layout.Inset{Bottom: unit.Dp(12), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(th, "Searching for local devices"+dots)
						lbl.Color = themes.ColorTextMuted
						lbl.TextSize = unit.Sp(13)
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						countStr := fmt.Sprintf("%d connected", len(devices))
						if len(devices) == 1 {
							countStr = "1 connected"
						}
						lbl := material.Caption(th, countStr)
						lbl.Color = themes.ColorBrown
						return lbl.Layout(gtx)
					}),
				)
			})
		}),

		// Scrollable Devices List or Empty State
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if len(devices) == 0 {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body1(th, "No devices found yet")
								lbl.Color = themes.ColorText
								lbl.Alignment = text.Middle
								return lbl.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(th, "Ensure ClipSync is running on your PC or other phones on the same Wi-Fi.")
								lbl.Color = themes.ColorTextMuted
								lbl.Alignment = text.Middle
								return lbl.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, connectBtn, "Connect Manually by IP")
								btn.Background = themes.ColorSurface
								btn.Color = themes.ColorBrown
								btn.TextSize = unit.Sp(13)
								btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}
								return btn.Layout(gtx)
							}),
						)
					})
				})
			}

			return list.Layout(gtx, len(devices), func(gtx layout.Context, index int) layout.Dimensions {
				return deviceCard(gtx, th, devices[index])
			})
		}),
	)
}

func deviceCard(gtx layout.Context, th *material.Theme, dev DeviceViewItem) layout.Dimensions {
	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 10, themes.ColorSurface, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(14)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								name := material.Body1(th, dev.Name)
								name.Color = themes.ColorText
								name.TextSize = unit.Sp(15)
								return name.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								ip := material.Caption(th, fmt.Sprintf("IP: %s", dev.IP))
								ip.Color = themes.ColorTextMuted
								ip.TextSize = unit.Sp(12)
								return ip.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						tag := material.Caption(th, "● Online")
						tag.Color = themes.ColorGreen
						tag.TextSize = unit.Sp(12)
						return tag.Layout(gtx)
					}),
				)
			})
		})
	})
}
