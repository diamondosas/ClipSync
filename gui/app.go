package gui

import (
	"clipsync/gui/components"
	"clipsync/gui/themes"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2"
)

func StartGUI() {
	// Create the main application
	a := app.New()

	// Apply our custom theme from gui/themes/theme.go
	a.Settings().SetTheme(&themes.MyTheme{})

	// Set up the main window defined in gui/components/window.go
	w := components.SetupWindow(a)
	if desk, ok := a.(desktop.App); ok {
		// 2. Create the menu items
		menu := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show Window", func() {
				w.Show()
			}),
		)
		// 3. Set the menu to the tray
		desk.SetSystemTrayMenu(menu)
	}



	w.SetCloseIntercept(func() {
		 w.Hide()
})

	// Run the application (this blocks the main thread)
	w.ShowAndRun()
}
