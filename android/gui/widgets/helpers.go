package widgets

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ColorBox fills a region with a solid color and renders the inner widget.
func ColorBox(gtx layout.Context, c color.NRGBA, inner layout.Widget) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, c, clip.Rect{Max: gtx.Constraints.Min}.Op())
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(inner),
	)
}

// RoundedBox fills a region with a solid color and rounded corners.
func RoundedBox(gtx layout.Context, radius int, c color.NRGBA, inner layout.Widget) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			r := clip.RRect{
				Rect: image.Rectangle{Max: gtx.Constraints.Min},
				NW:   radius, NE: radius, SW: radius, SE: radius,
			}
			paint.FillShape(gtx.Ops, c, r.Op(gtx.Ops))
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(inner),
	)
}

// TouchTarget ensures an element occupies at least minHeight for mobile ergonomics.
func TouchTarget(gtx layout.Context, minHeight unit.Dp, inner layout.Widget) layout.Dimensions {
	minPx := gtx.Dp(minHeight)
	if gtx.Constraints.Min.Y < minPx {
		gtx.Constraints.Min.Y = minPx
	}
	return inner(gtx)
}
