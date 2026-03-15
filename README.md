# Project Overview
**ClipSync** is a LAN-based clipboard-sharing application that synchronizes clipboards across devices on the same network.  
Architecture: peer-to-peer discovery via **mDNS** (service discovery) and **TCP** (data transfer).

---

# Technical Architecture Assessment

| Component        | Assessment |
|------------------|------------|
| **Discovery**    | **mDNS** (via [zeroconf](https://github.com/grandcat/zeroconf)) is ideal for LAN—zero manual IP entry, seamless UX. |
| **Communication**| **TCP** using net ensures reliable delivery, critical for clipboard integrity. Must handle text, images, files, etc. |
| **Cross-platform GUI** | "Gio UI" gives native-looking UIs on Windows, Linux, macOS from a single Go codebase. |

---

I am using this as a repo and the link is https://github.com/DiamondOsas/ClipSync.git




You just need to install the X11 development package for your Linux distribution before building your Go program.

Ubuntu / Debian / Pop!_OS / Linux Mint
sudo apt update
sudo apt install libx11-dev

That installs the header files including X11/Xlib.h.

Fedora / CentOS / RHEL
sudo dnf install libX11-devel
Arch Linux / Manjaro
sudo pacman -S libx11