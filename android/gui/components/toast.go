package components

import (
	"clipsync-android/gui/themes"
	"clipsync-android/gui/widgets"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Toast renders a floating notification with smooth fade-in and fade-out animation.
func Toast(gtx layout.Context, th *material.Theme, message string, startTime time.Time) layout.Dimensions {
	if message == "" || startTime.IsZero() {
		return layout.Dimensions{}
	}

	elapsed := gtx.Now.Sub(startTime)
	totalDuration := 1800 * time.Millisecond
	fadeIn := 200 * time.Millisecond
	fadeOut := 400 * time.Millisecond

	if elapsed >= totalDuration {
		return layout.Dimensions{}
	}

	gtx.Execute(op.InvalidateCmd{})

	var alpha uint8 = 255
	if elapsed < fadeIn {
		alpha = uint8((float32(elapsed) / float32(fadeIn)) * 255)
	} else if elapsed > totalDuration-fadeOut {
		remaining := totalDuration - elapsed
		alpha = uint8((float32(remaining) / float32(fadeOut)) * 255)
	}

	bgColor := color.NRGBA{R: themes.ColorSurface.R, G: themes.ColorSurface.G, B: themes.ColorSurface.B, A: alpha}
	textColor := color.NRGBA{R: themes.ColorText.R, G: themes.ColorText.G, B: themes.ColorText.B, A: alpha}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 8, bgColor, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, message)
				lbl.Color = textColor
				lbl.TextSize = unit.Sp(13)
				return lbl.Layout(gtx)
			})
		})
	})
}
