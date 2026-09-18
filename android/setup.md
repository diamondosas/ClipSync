To build and run a Gio application on Android from Linux, you need Go, OpenJDK, the Android SDK/NDK, and Gio's packaging tool `gogio`.

---

### 1. Install Java and Android Tools

Gio requires OpenJDK and the Android command-line utilities (including `adb` and `sdkmanager`).

On Debian/Ubuntu-based distributions:

```bash
sudo apt update
sudo apt install -y openjdk-17-jdk android-sdk adb

```

On Arch/Manjaro:

```bash
sudo pacman -S jdk17-openjdk android-tools

```

**Verification:** Run `adb version`. It should print the installed Android Debug Bridge version.

---

### 2. Configure Android SDK & Install NDK

Set up your `ANDROID_HOME` path and install the required Android platform headers and NDK.

Add the environment paths to your shell profile (e.g., `~/.bashrc` or `~/.zshrc`):

```bash
export ANDROID_HOME=$HOME/Android/Sdk
export PATH=$PATH:$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools

```

Apply the changes and install the build tools and NDK:

```bash
source ~/.bashrc
mkdir -p "$ANDROID_HOME"

# Accept SDK licenses and install NDK + platform tools
sdkmanager --licenses
sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0" "ndk;26.1.10909125"

```

**Verification:** Run `sdkmanager --list_installed`. You should see `ndk`, `platforms;android-34`, and `build-tools` in the list.

---

### 3. Install the Gio Build Tool (`gogio`)

`gogio` handles compiling the Go code and bundling it into a signed, installable `.apk` file.

```bash
go install gioui.org/cmd/gogio@latest

```

Ensure `$GOPATH/bin` or `$HOME/go/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin

```

**Verification:** Run `which gogio`. The terminal should return the path to the installed binary.

---

### 4. Create the Hello World Project

Create a new directory and initialize the module:

```bash
mkdir gio-hello && cd gio-hello
go mod init gio-hello

```

Create `main.go`:

```go
package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		window := new(app.Window)
		if err := run(window); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			label := material.H3(theme, "Hello, Android!")
			label.Color = color.NRGBA{R: 20, G: 120, B: 220, A: 255}
			label.Alignment = text.Middle
			label.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

```

Fetch the dependencies:

```bash
go mod tidy

```

**Verification:** Test that the desktop version compiles and runs on your Linux machine first by running `go run .`. A window displaying "Hello, Android!" should pop up.

---

### 5. Build and Install the APK

1. Enable **Developer Options** and **USB Debugging** on your Android phone.
2. Connect the phone to your Linux PC with a USB cable.
3. Build the APK using `gogio`:

```bash
gogio -target android -appid com.example.giohello -o hello.apk .

```

4. Install the APK to your connected phone:

```bash
adb install -r hello.apk

```

**Verification:** Run `adb devices` to verify your phone is detected. Once `adb install` outputs `Success`, open the **Hello Gio** app from your phone's app drawer.