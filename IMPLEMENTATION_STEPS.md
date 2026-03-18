# Implementation Steps for ClipSync GUI

This document provides the exact steps and code snippets needed to set up your Fyne application.

## 1. Update `main.go`
Open `main.go` and ensure the `StartGUI()` function is called at the end of the `main()` function. This starts the GUI once your background services are ready.

**What to do:**
- Import `"clipsync/gui"` at the top.
- Add `gui.StartGUI()` as  the last line in `main()`.

## 2. Implement `gui/app.go`
This file is the controller. It initializes the app, sets the theme, and starts the window.

```go
package gui

import (
	"clipsync/gui/components"
	"clipsync/gui/themes"

	"fyne.io/fyne/v2/app"
)

func StartGUI() {
	// Create the main application
	a := app.New()

	// Apply our custom theme from gui/themes/theme.go
	a.Settings().SetTheme(&themes.MyTheme{})

	// Set up the main window defined in gui/components/window.go
	w := components.SetupWindow(a)

	// Run the application (this blocks the main thread)
	w.ShowAndRun()
}
```

## 3. Implement `gui/themes/theme.go`
Create this file to define your custom look. It mimics the "learn" project theme but is simpler.

```go
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
	}
	return theme.DarkTheme().Color(name, variant)
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DarkTheme().Font(style) }
func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DarkTheme().Icon(name) }
func (m MyTheme) Size(name fyne.ThemeSizeName) float32 { return theme.DarkTheme().Size(name) }
```

## 4. Implement `gui/components/window.go`
This file builds the UI elements.

```go
package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SetupWindow(a fyne.App) fyne.Window {
	// Create a new window with a title
	w := a.NewWindow("🚀 ClipSync")

	// Set a reasonable initial size
	w.Resize(fyne.NewSize(400, 300))

	// Create a simple layout with a label and a separator
	content := container.NewVBox(
		widget.NewLabelWithStyle("ClipSync is Running", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("Syncing clipboard between your devices."),
		widget.NewButton("Open Settings", func() {
			// Placeholder for future logic
		}),
	)

	// Apply padding to the layout
	paddedContent := container.NewPadded(content)

	w.SetContent(paddedContent)
	return w
}
```

## 5. Final Checklist
1. Ensure all packages are correctly imported.
2. Run `go mod tidy` in your terminal to fetch any missing dependencies.
3. Start the app by running `go run .` in the root directory.
