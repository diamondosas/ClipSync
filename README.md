# ClipSync

<p align="center">
  <img src="asset/logo.jpg" alt="ClipSync Logo" width="128">
</p>

Stop emailing yourself links. Stop Slacking yourself snippets. Stop the friction.

ClipSync is the invisible peer-to-peer layer that synchronizes your clipboard across every device on your local network. It is built for those who value flow state over file transfers. No accounts. No cloud. No latency.

### The Manifesto

We believe that data should move as fast as you do. In a world of over-engineered cloud solutions, ClipSync returns to the basics: raw performance and local reliability. If your devices are on the same network, they should share the same clipboard. Period.

### Why ClipSync

*   **Zero Configuration**: mDNS service discovery means your devices find each other instantly. No manual IP entry. No handshake. Just sync.
*   **True Privacy**: Your clipboard data never leaves your local network. It moves directly from device A to device B.
*   **Industrial Reliability**: Built on Go with raw TCP transport to ensure every byte arrives exactly as it was copied.
*   **Native Performance**: Light on resources, heavy on utility. It sits in the background and does its job.

### Installation

Clone the repository and run the engine.

```bash
git clone https://github.com/DiamondOsas/ClipSync.git
cd ClipSync
go run main.go
```

### System Requirements

For those running Linux, the following development packages are required to interface with the X11 clipboard system.

**Debian / Ubuntu / Pop!_OS**
```bash
sudo apt update && sudo apt install libx11-dev
```

**Fedora / CentOS / RHEL**
```bash
sudo dnf install libX11-devel
```

**Arch Linux / Manjaro**
```bash
sudo pacman -S libx11
```

### Development

ClipSync is currently in its early stages of development. We are building the future of local-first productivity. 

[GitHub Repository](https://github.com/DiamondOsas/ClipSync.git)
