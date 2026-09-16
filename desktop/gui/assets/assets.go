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

// PinImageOp is the pre-decoded pushpin image operation for Gio.
var PinImageOp paint.ImageOp

func init() {
	img, _, err := image.Decode(bytes.NewReader(PinIconBytes))
	if err != nil {
		log.Printf("[Assets] Failed to decode pin.png: %v", err)
		return
	}
	PinImageOp = paint.NewImageOp(img)
}
