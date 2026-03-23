# ClipSync GUI Architecture

This document explains the structural layout of the new GUI implemented using the Gio UI library.

## Directory Structure
The `gui` package has been modularized to ensure separation of concerns and maintainability.

* `gui/app.go`: The main entry point. Sets up the window, main event loop, and ties everything together.
* `gui/state.go`: Defines `AppState` which manages all interactive state (tabs, dialog visibility, mock lists data) centrally.
* `gui/themes/`: Holds all static visual configurations (colors, future theme specs).
* `gui/widgets/`: Holds low-level, highly reusable layout helpers like `ColorBox` and `RoundedBox`.
* `gui/components/`: Reusable, slightly higher-level UI pieces (e.g., `Header`, `FooterTabs`, `HelpDialog`).
* `gui/pages/`: Specific layouts and views corresponding to entire app sections (e.g., `DevicesPage`, `ClipboardPage`).

## How it works (Data Flow)

1. **State Initialization**: The `AppState` holds your mock data (devices and clipboard history) along with clickable button states. It is initialized in `gui/app.go`.
2. **Event Handling**: On every UI frame, `state.Update(gtx)` is called. It iterates through tab buttons and the help buttons, updating active states or visibility flags (`ShowHelp`).
3. **Layout Rendering**: After processing clicks, we render the screen by executing `layoutMain(...)`. 
   - `layoutMain` wraps the application in the dark background using a layout wrapper (`widgets.ColorBox`).
   - It organizes the UI as a Vertical Flex stack: Rigid Header -> Flexed Body (occupying remaining space) -> Rigid Footer.
   - It also wraps the entire visual layout within `components.HelpDialog`. If `ShowHelp` is true, the dialog overlay draws on top of the main UI; otherwise, it just draws the main UI seamlessly.

## Adding Features later
* **New Page**: Create a new layout function in `gui/pages/`, add a corresponding tab button to the array in `gui/state.go`, and map the state in the Body flex block inside `gui/app.go`.
* **Real Data**: When you connect the network scanner and clipboard sync modules, simply update the arrays (`Devices` and `History`) stored in `AppState` inside a goroutine and trigger an `Invalidate()` to redraw the window.