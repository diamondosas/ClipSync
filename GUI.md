# ClipSync GUI Plan

This document outlines the design and structure of the ClipSync application interface. The goal is to create a simple, clean, and functional UI using the Gio library.

## General Design
- **Primary Color:** Electric Cyan (`#00C2FF`) for highlights and titles.
- **Theme:** Dark theme with grey accents for secondary text.
- **Layout:** A vertical layout with a persistent top header, a middle content area, and a bottom navigation bar.

---

## 1. Main Header (Persistent)
The top of the application will always show:
- **Title:** "ClipSync" in **Electric Cyan** bold text.
- **Help Button:** A small "?" or "Help" icon in the top right corner.
  - **Action:** Clicking this opens a small pop-up (dialog).
  - **Help Content:** Instructions on "How to connect" (e.g., "Ensure both devices are on the same Wi-Fi").
  - **Close Button:** A simple "OK" button to return to the main screen.

---

## 2. Navigation (Bottom Tabs)
The bottom of the screen will feature two tabs for switching between views. Each tab will have an **Icon** and a **Label**:
1. **Connection Info:** Icon + "Devices" (Default Page).
2. **Clipboard:** Icon + "Clipboard".

---

## 3. Page 1: Connection Info (Devices)
This page helps users see what devices are currently synced.
- **Status Text:** Below the header, it will display "Looking for devices...".
- **Device List:** A scrollable list of found devices.
- **Device Card Layout:**
  - **Device Name:** Bold text (e.g., "My Laptop").
  - **IP Address:** Smaller text below the name (e.g., "192.168.1.5").
  - *Note: This list is for viewing only; clicking a device does nothing for now.*

---

## 4. Page 2: Clipboard History
This page shows a history of items copied across devices.
- **Title:** "Clipboard" text at the top of the content area.
- **Clipboard List:** A scrollable list of "cards" representing saved snippets.
- **Clipboard Card Layout:**
  - **Snippet:** A preview of the text content.
  - **Interaction:** Clicking the card will simulate "Copy to Clipboard" (mock action for now).
  - **Future Proofing:** The code will be structured to easily add a dedicated "Copy" button later.

---

## 5. Technical Implementation Details (Internal)
- **Scrollability:** Both the Device List and Clipboard History will be contained within scrollable areas so they don't overflow the window.
- **Mock Data:** The initial implementation will use hardcoded lists for devices and clipboard items to demonstrate the UI.
- **Window Size:** The app will be designed for a "Final App" feel, likely a compact, tall window suitable for a utility tool.
