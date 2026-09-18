# Building ClipSync for Android

This guide explains how to compile the ClipSync Android application with full **Accessibility-powered Automatic Background Clipboard Sync** using `gogio`.

---

## 1. Prerequisites

Ensure you have the following installed:
1. **Go (1.22+)**
2. **Android SDK & NDK** (`ANDROID_HOME` configured in your shell)
3. **`gogio` tool**:
   ```bash
   go install gioui.org/cmd/gogio@latest
   ```

---

## 2. Compile the APK

```bash
cd android

# Step 1: Build the Gio Core AAR library
mkdir -p app/libs
gogio -target android -buildmode archive -o app/libs/clipsync.aar .

# Step 2: Assemble the final APK with Accessibility Service
./gradlew assembleRelease
```

The compiled APK will be generated at:
`android/app/build/outputs/apk/release/app-release.apk`

---

## 3. Install on Android Phone

1. Connect your Android phone via USB with **USB Debugging** enabled.
2. Verify connection:
   ```bash
   adb devices
   ```
3. Install the APK:
   ```bash
   adb install -r clipsync.apk
   ```

---

## 4. Enable 100% Automatic Background Sync (One-Time Step)

Due to Android 10+ clipboard security restrictions:
1. Open phone **Settings**.
2. Go to **Accessibility** (or *Accessibility > Installed Services* on Samsung/Xiaomi/Pixel).
3. Find **ClipSync Automatic Clipboard Sync**.
4. Toggle it **ON** and tap **Allow**.

---

## 5. How It Works

```
┌───────────────────────────────────────────────┐
│ User Copies Text Anywhere in Any Android App   │
└───────────────────────┬───────────────────────┘
                        │
                        ▼
┌───────────────────────────────────────────────┐
│ ClipSyncAccessibilityService (Java)           │
│ Detects clipboard change event                │
└───────────────────────┬───────────────────────┘
                        │ UDP (127.0.0.1:9998)
                       clipsync ▼
┌───────────────────────────────────────────────┐
│ Go LocalBridge & ClipSync Engine              │
│ Deduplicates, saves to history                │
└───────────────────────┬───────────────────────┘
                        │ Encrypted KCP/UDP (Port 9999)
                        ▼
┌───────────────────────────────────────────────┐
│ Desktop PC (Linux, Windows, macOS)            │
│ Instantly receives & syncs clipboard!         │
└───────────────────────────────────────────────┘
```
