package themes

import "image/color"

var (
	// Root background: Pure Black
	ColorBg = color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	// Primary Accent: Warm Amber Caramel
	ColorBrown  = color.NRGBA{R: 198, G: 118, B: 40, A: 255}
	ColorAccent = ColorBrown
	ColorCyan   = ColorBrown // Compatibility alias

	// Primary Text: Bright Warm Ivory Cream
	ColorCream = color.NRGBA{R: 255, G: 246, B: 236, A: 255}
	ColorText  = ColorCream

	// Secondary Text: Soft Warm Almond
	ColorTextMuted = color.NRGBA{R: 218, G: 201, B: 185, A: 255}

	// Surface & Card Backgrounds: Warm Rich Chocolate Brown
	ColorSurface       = color.NRGBA{R: 54, G: 34, B: 20, A: 255}
	ColorSurfacePinned = color.NRGBA{R: 90, G: 52, B: 22, A: 255}

	// Destructive Action: Warm Rustic Terracotta Red
	ColorRed = color.NRGBA{R: 220, G: 75, B: 60, A: 255}

	// Status Indicator: Emerald Green
	ColorGreen = color.NRGBA{R: 56, G: 142, B: 60, A: 255}
)
