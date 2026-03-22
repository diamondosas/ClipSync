# GUI Update: Switching to Gio

I have replaced the **Fyne** GUI framework with **Gio** (as seen in `learn.go`). Gio is a more lightweight, immediate-mode GUI framework that allows for very minimal code.

### Changes Made:
1.  **Framework Swap**: Moved from `fyne.io/fyne/v2` to `gioui.org`.
2.  **Simplified Structure**: Consolidated the GUI logic into `gui/app.go`.
3.  **Minimal `main.go`**: The entry point remains a single call to `gui.StartGUI()`.

### How it works:
- `gui.StartGUI()`: Starts the Gio application loop.
- `gui.run()`: Handles the window events and renders a simple "ClipSync is Running" message.
- **Immediate Mode**: Unlike Fyne, Gio redraws the UI every frame (or when events happen), giving you full control over the layout.

### Next Steps:
You can now easily add more widgets by following the patterns in `learn/learn.go`.
