// ai-generated
package themes

import "image/color"

var (
	// Root background: Pure Black
	ColorBg = color.NRGBA{R: 0, G: 0, B: 0, A: 255} // #000000 Pure Black Canvas

	// Primary Accent: Vibrant Amber Caramel (#C67628 / enhanced #B26925)
	ColorBrown  = color.NRGBA{R: 198, G: 118, B: 40, A: 255} // #C67628 Glowing Warm Amber
	ColorCyan   = ColorBrown                                 // Backward compatibility alias
	ColorAccent = ColorBrown

	// Primary Text: Bright Warm Ivory Cream
	// Crisp, radiant, high-contrast text that never looks dull on dark backgrounds
	ColorCream = color.NRGBA{R: 255, G: 246, B: 236, A: 255} // #FFF6EC Luminous Warm Ivory
	ColorText  = ColorCream

	// Secondary Text: Soft Warm Almond / Latte (#DAC9B9)
	// Promoted to secondary text so subtitles, placeholders, and buttons are bright and clear
	ColorTextMuted = color.NRGBA{R: 218, G: 201, B: 185, A: 255} // #DAC9B9 Almond Latte Cream

	// Surface & Card Backgrounds: Warm Rich Chocolate Brown (distinct from pure black)
	ColorSurface       = color.NRGBA{R: 54, G: 34, B: 20, A: 255} // #362214 Rich Warm Chocolate Brown
	ColorSurfacePinned = color.NRGBA{R: 90, G: 52, B: 22, A: 255} // #5A3416 Warm Bronze Glow for Pinned Cards

	// Destructive Action: Warm Rustic Terracotta Red
	ColorRed = color.NRGBA{R: 220, G: 75, B: 60, A: 255} // #DC4B3C Warm Crimson Red
)
