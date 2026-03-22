package widgets

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// ColorBox fills a region with a solid color and draws the inner widget.
func ColorBox(gtx layout.Context, c color.NRGBA, inner layout.Widget) layout.Dimensions {
	m := op.Record(gtx.Ops)
	dims := inner(gtx)
	call := m.Stop()
	
	// Draw the background color
	defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
	
	// Add the widget content on top
	call.Add(gtx.Ops)
	return dims
}

// RoundedBox fills a region with a solid color and rounded corners.
func RoundedBox(gtx layout.Context, radius int, c color.NRGBA, inner layout.Widget) layout.Dimensions {
	m := op.Record(gtx.Ops)
	dims := inner(gtx)
	call := m.Stop()
	
	// Draw the rounded background color
	defer clip.RRect{
		Rect: image.Rect(0, 0, dims.Size.X, dims.Size.Y),
		NW:   radius, NE: radius, SW: radius, SE: radius,
	}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
	
	// Add the widget content on top
	call.Add(gtx.Ops)
	return dims
}
