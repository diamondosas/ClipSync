// ai-generated
package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"gioui.org/op/paint"
)

//go:embed pin.png
var PinIconBytes []byte

//go:embed cancel.png
var CancelIconBytes []byte

// PinImageOp is the pre-decoded pushpin image operation for Gio.
var PinImageOp paint.ImageOp

// CancelImageOp is the pre-decoded cancel/delete image operation for Gio.
var CancelImageOp paint.ImageOp

func init() {
	if img, _, err := image.Decode(bytes.NewReader(PinIconBytes)); err != nil {
		log.Printf("[Assets] Failed to decode pin.png: %v", err)
	} else {
		PinImageOp = paint.NewImageOp(img)
	}

	if img, _, err := image.Decode(bytes.NewReader(CancelIconBytes)); err != nil {
		log.Printf("[Assets] Failed to decode cancel.png: %v", err)
	} else {
		CancelImageOp = paint.NewImageOp(img)
	}
}
