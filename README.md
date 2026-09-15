<h2 align="center">
  <b>Copy on your Android. Paste on your Desktop.</b><br>
  <sub>No account. No friction. Same clipboard content on all devices connected.</sub>
</h2>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Built with Go">
  <img src="https://img.shields.io/badge/Network-LAN%20Only%2C%20No%20Cloud-green?style=for-the-badge" alt="Local network only">
  <a href="https://diamondosas.github.io/clipsync/">
    <img src="https://img.shields.io/badge/Download-Now-blue?style=for-the-badge&logo=appveyor" alt="Download ClipSync">
  </a>
</p>

<!-- ---

**ClipSync** is an open-source clipboard tool that syncs across every device on your local network — Windows, macOS, Linux — instantly. Built in Go. Runs silent. Uses barely any CPU,  RAM. No servers.

--- -->

<h3 align="center">Screenshots</h3>

<p align="center">
  <img src="assets/screenshots/devices.png" alt="Devices View" width="380" style="border-radius: 8px; margin: 6px;" />
  <img src="assets/screenshots/clipboard.png" alt="Clipboard View" width="380" style="border-radius: 8px; margin: 6px;" />
</p>


## Download ClipSync

Get the latest binary for your OS — run it, done.

**[Download ClipSync — Windows, macOS, Linux](https://diamondosas.github.io/clipsync/)**

or

**[Go to Realeases](https://github.com/diamondosas/clipsync/releases)**

---

---

## Features

- Automatic device discovery
- Lightweight background process (~15 MB RAM, ~0.1% CPU)
- Fast clipboard sync across devices
- Local-network only—your data stays private
- Works on Windows, macOS, and Linux
- Supports both X11 and Wayland


---

---

## Linux Setup — One Extra Step

On Linux, ClipSync needs access to the X11 or Wayland clipboard layer. Install the right libraries for your distro first:


**Debian / Ubuntu / Pop!_OS / Mint**
```bash
sudo apt install libx11-dev libwayland-dev libxkbcommon-dev libvulkan-dev
```

**Fedora / CentOS / RHEL / AlmaLinux**
```bash
sudo dnf install libX11-devel libwayland-dev libxkbcommon-dev vulkan-headers
```

**Arch / Manjaro / EndeavourOS**w
```bash
sudo pacman -S libx11 libwayland-dev libxkbcommon-dev libvulkan-dev
```

Display Server Libraies _check using :_ `echo $XDG_SESSION_TYPE`

**Wayland** 

```bash
sudo apt install wl-clipboard
```

**X11** 


---
---

## ⭐ Star History

If ClipSync saves you time, star it. Helps other people find it.

<!-- <picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=diamondosas/clipsync&type=date&theme=dark&legend=top-left" />
  <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=diamondosas/clipsync&type=date&legend=top-left" />
  <img alt="ClipSync GitHub star history" src="https://api.star-history.com/chart?repos=diamondosas/clipsync&type=date&legend=top-left" />
</picture> -->

---
