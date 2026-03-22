# ClipSync GUI Implementation Plan

This guide explains how we will build the ClipSync interface using the **Gio** library, inspired by the examples in `learn.go`.

## 1. Project Structure
We will organize our GUI code in the `gui/` directory.
- `gui/app.go`: The main entry point for the GUI.
- `gui/components/`: Reusable UI pieces (buttons, cards, etc.).

---

## 2. Setting Up the Theme
We will use a dark theme with **Electric Cyan** (`#00C2FF`) as our primary accent color.

```go
var (
    colorBg      = color.NRGBA{R: 18, G: 18, B: 18, A: 255}    // Dark Background
    colorCyan    = color.NRGBA{R: 0, G: 194, B: 255, A: 255}   // Electric Cyan
    colorSurface = color.NRGBA{R: 30, G: 30, B: 30, A: 255}    // Card Background
    colorText    = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White Text
)
```

---

## 3. Main Application State
Just like in the "learn" example, we need a struct to keep track of what's happening in our app (like which tab is open).

```go
type AppState struct {
    th        *material.Theme
    activeTab int
    tabBtns   [2]widget.Clickable
    
    // Page 1: Devices
    deviceList widget.List
    devices    []Device // Mock data for now
    
    // Page 2: Clipboard
    clipList   widget.List
    history    []string // Mock history
}
```

---

## 4. Building the Layout
We will use a `layout.Flex` to stack three main parts vertically:
1. **Header:** Title and Help button.
2. **Content:** The active page (Devices or Clipboard).
3. **Footer:** The navigation tabs.

```go
func (s *AppState) Layout(gtx layout.Context) layout.Dimensions {
    return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
        layout.Rigid(s.layoutHeader),   // Top
        layout.Flexed(1, s.layoutBody), // Middle (grows to fill space)
        layout.Rigid(s.layoutFooter),   // Bottom
    )
}
```

---

## 5. Components & Pages

### Header (Persistent)
The header will show the "ClipSync" title in Cyan.
```go
func (s *AppState) layoutHeader(gtx layout.Context) layout.Dimensions {
    return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        title := material.H5(s.th, "ClipSync")
        title.Color = colorCyan
        return title.Layout(gtx)
    })
}
```

### Device Cards
For the device list, we will create a simple card for each device.
```go
func (s *AppState) deviceCard(gtx layout.Context, name, ip string) layout.Dimensions {
    return colorBox(gtx, colorSurface, func(gtx layout.Context) layout.Dimensions {
        return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(material.Body1(s.th, name).Layout),
                layout.Rigid(material.Caption(s.th, ip).Layout),
            )
        })
    })
}
```

---

## 6. How to Run
Once implemented, you can run the GUI by calling:
```bash
go run main.go
```
(Assuming `main.go` initializes the GUI from the `gui/` package).

## Simple Summary for the User
1. **Define Colors:** Set up the Dark + Cyan look.
2. **State Management:** Create a struct to hold button clicks and lists.
3. **Flex Layout:** Stack the Header, Body, and Navigation.
4. **Lists:** Use `widget.List` to show scrollable devices and clipboard items.
