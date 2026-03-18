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
	paddedContent := container.New(container.NewPadded().Layout, content)

	w.SetContent(paddedContent)
	return w
}
