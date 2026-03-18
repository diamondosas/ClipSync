# Fyne GUI Implementation Guide

This guide explains how to set up your ClipSync GUI using the Fyne toolkit. We will break it down into four simple parts.

## Overview
1.  **Entry Point (`main.go`)**: Starts the background services and the GUI.
2.  **App Configuration (`gui/app.go`)**: Initializes the Fyne application and sets the theme.
3.  **Theme Definition (`gui/themes/theme.go`)**: Customizes colors and sizes.
4.  **Window Setup (`gui/components/window.go`)**: Configures the main window content and features.

---

## Step 1: Update `main.go`
In your `main.go`, you need to call `StartGUI()` at the end of the `main` function. Note that Fyne must run on the **main thread**, so we call it after starting our background tasks.

**What to do:**
- Import the `gui` package (e.g., `clipsync/gui`).
- Call `gui.StartGUI()` as the last line of `main`.

---

## Step 2: Configure the App in `gui/app.go`
This file acts as the bridge. It creates the Fyne application and applies your custom theme.

**What to do:**
1. Define `StartGUI()`.
2. Create a new app using `app.New()`.
3. Apply the theme from `gui/themes`.
4. Call the window setup from `gui/components`.
5. Run the application using `w.ShowAndRun()`.

---

## Step 3: Define the Theme in `gui/themes/theme.go`
Inspired by `learn.go`, this file implements the `fyne.Theme` interface to give your app a unique look.

**What to do:**
- Create a `MyTheme` struct.
- Implement `Color`, `Font`, `Icon`, and `Size` methods.
- For a "simple" look, focus on `theme.ColorNamePrimary` and `theme.ColorNameBackground`.

---

## Step 4: Setup the Window in `gui/components/window.go`
This is where you decide what the user actually sees.

**What to do:**
- Create a function like `SetupWindow(a fyne.App)`.
- Use `a.NewWindow("ClipSync")` to create the window.
- Set the content using containers (like `container.NewVBox`).
- Add a simple label or button to verify it works.

---

## Technical Flow Diagram
`main.go` -> calls `gui.StartGUI()`
  `gui.StartGUI()` -> creates `app.New()`
    `app.New()` -> applies `theme.MyTheme`
      `app.New()` -> calls `components.SetupWindow()`
        `components.SetupWindow()` -> defines UI layout and returns the window
          `gui.StartGUI()` -> calls `window.ShowAndRun()`
