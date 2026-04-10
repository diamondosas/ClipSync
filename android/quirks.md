# Android Quirks & Best Practices Guide for Gio UI

## Critical Android Considerations

### 1. Safe Area Insets (Notches, Cutouts, System UI)

**The Problem:**
Android devices have various safe area boundaries:
- Status bar at top (usually 24-25dp)
- Navigation bar at bottom (varies: 48dp on some, 0 on gesture nav)
- Notches/punch-holes at top
- Display cutouts on edges

**Gio's Solution:**
```go
// Receive insets from FrameEvent
case system.FrameEvent:
    state.insets = e.Insets

// Apply insets to your root layout
inset := layout.Inset{
    Top:    unit.Dp(float32(state.insets.Top)),
    Bottom: unit.Dp(float32(state.insets.Bottom)),
    Left:   unit.Dp(float32(state.insets.Left)),
    Right:  unit.Dp(float32(state.insets.Right)),
}
inset.Layout(gtx, yourContent)
```

**Why It Matters:**
- Without insets, content can be hidden behind status/nav bars
- Notches can obscure UI elements in notched areas
- Users expect content to respect device boundaries

### 2. Touch Target Sizing (Material Design Spec)

**The Rule:**
Minimum touch target = 48dp × 48dp (or 44dp × 44dp minimum)

**Why 48dp?**
- Average finger width is ~10mm (40-50dp)
- Prevents mis-taps and accidental interactions
- Required for accessibility compliance

**Implementation:**
```go
// Good: Slider with 44dp+ height
gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(44))
sl := material.Slider(state.theme, &widget.Float{})
sl.Layout(gtx)

// Good: Button with adequate padding
btn := material.Button(state.theme, clickable, "Tap Me")
btn.Layout(gtx)  // Material button has built-in sizing

// Bad: Don't create small touch targets
layout.Dimensions{Size: image.Pt(20, 20)}
```

**Components to Check:**
- Buttons: Material buttons are auto-sized correctly
- Sliders: Explicitly set Min.Y ≥ 44dp
- Checkboxes/Radio: Material components handle this
- List items: Row height ≥ 48dp recommended

### 3. Keyboard & Input Handling

**Single-line vs Multiline:**
```go
// For single-line text (triggers standard keyboard)
textField.SingleLine = true

// For passwords (triggers hidden input keyboard)
textField.SingleLine = true
// Note: Gio doesn't have built-in password masking yet

// For multiline (disables submit-on-Enter)
editor.SingleLine = false
editor.Submit = false  // Prevents accidental submission
```

**Problem:**
On Android, pressing Enter in a multiline field closes the keyboard and submits
(which may not be desired). Set `Submit = false` to prevent this.

**Keyboard Types:**
Android shows different keyboards based on:
- Input content type (numbers, email, URL, etc.)
- Currently, Gio doesn't provide fine-grained control over keyboard type
- Future versions may add `InputType` or similar

**Solution for Number Input:**
```go
// Until Gio adds numeric input widget, simulate with editor:
numberField := widget.Editor{}
numberField.SingleLine = true
// Users will need to manually filter input or use runtime validation
```

### 4. Back Button Handling

**Android Back Press:**
```go
case key.Event:
    if e.Name == key.NameBack {
        // Clean up resources if needed
        return nil  // Exit app gracefully
    }
```

**Common Pattern:**
- Single back press = exit app
- Multiple screens? Use a navigation stack
- Modal dialogs? Dismiss on back press (not yet fully supported in base Gio)

### 5. Rotation & Screen Configuration Changes

**The Challenge:**
Android destroys and recreates the activity on rotation, but Gio's app.Window
persists. This means:
- Insets change (swap width/height constraints)
- Layout needs to reflow
- State should survive rotation

**Safe Practice:**
```go
// Store everything in AppState struct
type AppState struct {
    textField widget.Editor
    sliderValue float32
    selectedTab int
    insets system.Insets  // Will update on rotation
}

// Don't store context, operations, or GTX objects in state
// Only store widget state and model data
```

**What Gio Handles Automatically:**
- Calling FrameEvent with new constraints
- Updating insets for new orientation
- Widget state preservation (Editor, List scroll position, etc.)

**You Must Handle:**
- Re-reading insets in each FrameEvent
- Adjusting layout constraints if needed
- Re-creating cached theme if needed (usually not needed)

### 6. Memory & Performance on Mobile

**Limited Resources:**
- RAM: 2GB-8GB typical, but other apps competing
- Heap: Garbage collection pauses matter more
- Battery: Rendering is expensive

**Best Practices:**
```go
// ✓ Good: Create theme once
theme := material.NewTheme()

// ✗ Bad: Creating theme every frame
for e := w.Event() {
    theme := material.NewTheme()  // Allocation on every frame!
}

// ✓ Good: Use material.List for large datasets
list := widget.List{}
material.List(theme, &list).Layout(gtx, len(items), ...)

// ✗ Bad: Creating widgets for every item
for i := range items {
    widget.Button{}  // Allocation + GC pressure
}
```

**Scroll Performance:**
- Gio's `material.List` only renders visible items (viewport culling)
- This scales to thousands of items without stuttering
- Don't bypass `material.List` for custom scrolling

### 7. Density Independence & Scaling

**Device Density (DPI):**
- ldpi: 120 DPI
- mdpi: 160 DPI (160DPI = 1 "dp")
- hdpi: 240 DPI
- xhdpi: 320 DPI
- xxhdpi: 480 DPI
- xxxhdpi: 640 DPI

**Always Use dp (Device-Independent Pixels):**
```go
// ✓ Good: Density-independent
margin := unit.Dp(16)
fontSize := unit.Dp(14)

// ✗ Bad: Pixel-dependent (wrong on different densities)
margin := 16  // Assumed to be pixels!
fontSize := 14
```

**Gio's Density Handling:**
```go
// Gio converts dp to px internally
scale := gtx.Metric.PxPerDp

// Use unit.Dp(value) for all dimensions
unit.Dp(16) // Always use this
```

### 8. Color & Theme for Android

**Dark Mode Support:**
```go
// Gio doesn't auto-detect dark mode (yet)
// Manually check or provide user preference
theme := material.NewTheme()
if userPrefersDark {
    theme.Palette.Bg = color.NRGBA{R: 20, G: 20, B: 20, A: 255}
    theme.Palette.Fg = color.NRGBA{R: 245, G: 245, B: 245, A: 255}
}
```

**Contrast & Accessibility:**
- Foreground: #1e1e1e or #f5f5f5 (depending on theme)
- Background: #f0f0f5 or #141414
- Interactive: #3f51b5 (blue) with white text
- Ensure 4.5:1 contrast ratio minimum (WCAG AA)

**System Colors:**
Gio doesn't read system theme colors, so apps look consistent regardless of device theme.

### 9. List Optimization for Large Datasets

**Don't Do This:**
```go
// ✗ Bad: Creates all widgets upfront
var items []widget.Clickable
for i := 0; i < 10000; i++ {
    items = append(items, widget.Clickable{})
}
```

**Do This Instead:**
```go
// ✓ Good: Let material.List handle rendering
state.items := []string{/* 10,000 items */}
material.List(theme, &state.listScroll).Layout(
    gtx, len(state.items),
    func(gtx layout.Context, index int) layout.Dimensions {
        // Only called for visible items!
        return material.Body2(theme, state.items[index]).Layout(gtx)
    },
)
```

### 10. Gesture vs Button Navigation

**Android Navigation Buttons (Bottom):**
- Traditional: 3-button navigation (back, home, recents)
- Modern (Pie+): Gesture navigation (swipe from edges)

**Gio's Handling:**
- Back button: Sent as `key.NameBack` event
- Home/Recents: Handled by system, app backgrounded
- Swipe: Not available to app (handled by system)

**Best Practice:**
```go
case key.Event:
    if e.Name == key.NameBack {
        // Only handle back button, don't try to intercept gestures
        return nil
    }
```

## Advanced Android Topics

### Landscape vs Portrait Layouts

**Current Constraints in Gio:**
```go
case system.FrameEvent:
    if e.Size.X > e.Size.Y {
        // Landscape mode
        drawLandscapeLayout(gtx, state)
    } else {
        // Portrait mode
        drawPortraitLayout(gtx, state)
    }
```

### Respecting Gesture Navigation Gestures

**Problem:**
- System gesture nav uses edge swipes
- Apps shouldn't consume these swipes

**Solution:**
Gio doesn't expose edge swipe handling, which is correct. Don't fight the system.

### Handling Keyboard Visibility

**Current Limitation:**
Gio doesn't expose keyboard show/hide events directly. The keyboard appears/disappears
automatically when you focus text input.

**Workaround:**
Use insets to detect keyboard visibility:
```go
if e.Insets.Bottom > gtx.Dp(unit.Dp(300)) {
    // Keyboard likely shown (rough heuristic)
    scrollToShowFocusedField()
}
```

### Minimizing Battery Usage

**Good Practices:**
1. Use `material.List` for efficient rendering
2. Avoid continuous animations without need
3. Batch drawing operations (Gio does this automatically)
4. Don't create unnecessary goroutines
5. Cache theme and other expensive objects

**Frame Rate:**
- Gio automatically limits to 60fps or screen refresh rate
- No need to manually throttle
- Battery impact is minimized

## Testing Checklist

- [ ] Test on actual Android device (not just emulator)
- [ ] Test with notch/punch-hole present
- [ ] Test in both portrait and landscape
- [ ] Test with device back button
- [ ] Test with soft keyboard shown
- [ ] Test with text input fields
- [ ] Test scrolling large lists (1000+ items)
- [ ] Test with different screen sizes (4" - 6.7")
- [ ] Test with device density (hdpi, xhdpi, xxhdpi)
- [ ] Check memory usage (Android Studio Profiler)
- [ ] Test dark mode (if implemented)

## Common Gotchas

1. **Forgetting safe area insets** → Content hidden behind status bar
2. **Small touch targets** → App fails accessibility requirements
3. **Creating widgets in layout function** → Memory pressure & GC pauses
4. **Not using material.List** → Stuttering with large lists
5. **Using pixels instead of dp** → Wrong sizing on different screens
6. **Storing GTX in state** → Crashes or corrupted rendering
7. **Submitting on Enter in multiline** → Bad UX on mobile
8. **No back button handling** → App force-closes on back press
9. **Creating theme every frame** → Unnecessary allocations
10. **Assuming portrait orientation** → Rotation breaks layout

## Performance Profiling

### Using Android Studio Profiler
```bash
# Build with profiling enabled (gogio has flags)
gogio -target android -profile ...
```

### Key Metrics to Monitor
- **Frame time**: Should be < 16ms (60fps)
- **Memory**: Watch for memory growth
- **CPU**: Should be low when idle
- **Battery**: Monitor drain during use

## Example: Handling Complex State with Rotation

```go
type AppState struct {
    // Persistent state (survives rotation)
    textInput      widget.Editor
    scrollPosition int
    selectedItem   int
    
    // Temporary state (recalculated per frame)
    insets         system.Insets
    
    // Never store these:
    // - *layout.Context
    // - *op.Ops
    // - theme (although can cache)
}

// Rotation flow:
// 1. System rotates device
// 2. Window size changes
// 3. FrameEvent fired with new constraints
// 4. Insets updated
// 5. Layout reflows automatically
// 6. State preserved in AppState
```

## Resources

- [Android Developer Guide](https://developer.android.com/guide)
- [Material Design](https://material.io/design)
- [Gio Documentation](https://gioui.org/doc)
- [Accessibility Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [Android Performance Guide](https://developer.android.com/topic/performance)