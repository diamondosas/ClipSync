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
