// ai-generated
package components

import (
	"image/color"
	"time"

	"clipsync/gui/themes"
	"clipsync/gui/widgets"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Toast renders an animated floating notification that smoothly fades in and out.
func Toast(gtx layout.Context, th *material.Theme, message string, startTime time.Time) layout.Dimensions {
	if message == "" || startTime.IsZero() {
		return layout.Dimensions{}
	}

	elapsed := gtx.Now.Sub(startTime)
	totalDuration := 1500 * time.Millisecond
	fadeIn := 200 * time.Millisecond
	fadeOut := 400 * time.Millisecond

	if elapsed >= totalDuration {
		return layout.Dimensions{}
	}

	// Request next frame to maintain 60fps animation
	gtx.Execute(op.InvalidateCmd{})

	// Calculate opacity alpha (0 - 255)
	var alpha uint8 = 255
	if elapsed < fadeIn {
		alpha = uint8((float32(elapsed) / float32(fadeIn)) * 255)
	} else if elapsed > totalDuration-fadeOut {
		remaining := totalDuration - elapsed
		alpha = uint8((float32(remaining) / float32(fadeOut)) * 255)
	}

	bgColor := color.NRGBA{R: 28, G: 36, B: 44, A: alpha}
	textColor := color.NRGBA{R: themes.ColorText.R, G: themes.ColorText.G, B: themes.ColorText.B, A: alpha}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widgets.RoundedBox(gtx, 6, bgColor, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(14), Right: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(th, message)
				lbl.Color = textColor
				return lbl.Layout(gtx)
			})
		})
	})
}
