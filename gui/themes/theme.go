package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

// Satisfy the fyne.Theme interface
var _ fyne.Theme = (*MyTheme)(nil)

func (m MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		// Electric cyan from learn.go
		return color.NRGBA{R: 0x00, G: 0xC2, B: 0xFF, A: 0xFF}
	case theme.ColorNameBackground:
		// Dark grey background
		return color.NRGBA{R: 0x12, G: 0x12, B: 0x1A, A: 0xFF}
	case theme.ColorNameForeground:
		// Light grey foreground
		return color.NRGBA{R: 0xE0, G: 0xE0, B: 0xFF, A: 0xFF}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x1E, G: 0x1E, B: 0x2E, A: 0xFF}
	}
	return theme.DarkTheme().Color(name, variant)
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DarkTheme().Font(style) }
func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DarkTheme().Icon(name) }
func (m MyTheme) Size(name fyne.ThemeSizeName) float32 { return theme.DarkTheme().Size(name) }
