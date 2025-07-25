# 🏗️ Arch Linux Package Installation Guide

## 🚀 Quick Installation (Recommended)

### Method 1: Automated Installation Script

```bash
# Download and run the installation script
curl -sSL https://github.com/yourusername/hyprland-monitor-tui/raw/main/install.sh | bash

# Or if you have the source:
./install.sh
```

### Method 2: Manual Build & Install

```bash
# Install Go via mise (if not already available)
mise use -g go@latest

# Install optional dependencies (recommended)
sudo pacman -S wlr-randr xorg-xrandr

# Build the application
go build -o hyprland-monitor-tui .

# Install system-wide
sudo install -Dm755 hyprland-monitor-tui /usr/local/bin/hyprland-monitor-tui
sudo install -Dm644 hyprland-monitor-tui.desktop /usr/share/applications/hyprland-monitor-tui.desktop
```

## 📦 Package Creation for AUR

### Create Package Archive

```bash
# Create source archive
tar -czf hyprland-monitor-tui-1.0.0.tar.gz \
    main.go model.go monitor.go go.mod go.sum \
    README.md LICENSE hyprland-monitor-tui.desktop

# Generate checksums
sha256sum hyprland-monitor-tui-1.0.0.tar.gz
```

### Build with makepkg

```bash
# Build package
makepkg -si

# Or just build without installing
makepkg

# Install built package
sudo pacman -U hyprland-monitor-tui-1.0.0-1-x86_64.pkg.tar.zst
```

## 🔧 Arch Linux Compatibility

### ✅ Verified Components

- **Hyprland Integration**: Native `hyprctl` support
- **Wayland Support**: Uses `wlr-randr` for fallback detection
- **X11 Compatibility**: Falls back to `xrandr` when needed
- **Go Runtime**: Built with Go 1.21+ (available in Arch repos)
- **Dependencies**: All optional deps available in official repos

### 📋 System Requirements

```bash
# Build requirements (managed via mise)
mise use -g go@latest

# Runtime requirements (assumed to be installed)
pacman -S hyprland

# Optional but recommended
pacman -S wlr-randr xorg-xrandr

# For development
pacman -S git base-devel
```

### 🎯 Hyprland-Specific Features

The application is specifically designed for Hyprland and includes:

- **Native `hyprctl monitors` parsing**
- **Real-time scale application via `hyprctl keyword`**
- **Hyprland configuration file management**
- **Environment detection (`HYPRLAND_INSTANCE_SIGNATURE`)**
- **Graceful fallbacks for non-Hyprland environments**

## 🧪 Testing on Arch Linux

### In Hyprland Environment

```bash
# Normal operation with full functionality
hyprland-monitor-tui

# Debug mode for troubleshooting
hyprland-monitor-tui --debug
```

### Outside Hyprland (Testing)

```bash
# Demo mode with sample data
hyprland-monitor-tui --no-hyprland-check

# This will show realistic Framework 13 + external monitor data
# Perfect for UI testing and development
```

### Verification Commands

```bash
# Check monitor detection methods
hyprctl monitors  # Should work in Hyprland
wlr-randr         # Should work in any Wayland compositor
xrandr            # Should work in X11 environments

# Test scaling application (Hyprland only)
hyprctl keyword monitor eDP-1,preferred,auto,2.0
```

## 🔍 Troubleshooting

### Common Issues

1. **"hyprctl not found"**
   ```bash
   sudo pacman -S hyprland
   ```

2. **"No monitors detected"**
   ```bash
   # Install fallback tools
   sudo pacman -S wlr-randr xorg-xrandr
   ```

3. **"Permission denied"**
   ```bash
   # Ensure user is in correct groups
   groups $USER
   # Should include: wheel, video, input
   ```

### Debug Information

```bash
# Check Hyprland status
echo $HYPRLAND_INSTANCE_SIGNATURE

# Check available detection tools
which hyprctl wlr-randr xrandr

# Test monitor detection
hyprland-monitor-tui --debug --no-hyprland-check
```

## 🏷️ Package Information

- **Package Name**: `hyprland-monitor-tui`
- **Version**: `1.0.0`
- **Architecture**: `x86_64`, `aarch64`
- **License**: `MIT`
- **Dependencies**: `hyprland`
- **Optional Dependencies**: `wlr-randr`, `xorg-xrandr`
- **Build Dependencies**: `go`

## 🗑️ Uninstallation

```bash
# If installed via script
./uninstall.sh

# If installed via pacman
sudo pacman -R hyprland-monitor-tui

# Manual removal
sudo rm /usr/local/bin/hyprland-monitor-tui
sudo rm /usr/share/applications/hyprland-monitor-tui.desktop
sudo rm -rf /usr/share/doc/hyprland-monitor-tui/
```

---

**Ready to create the most beautiful monitor management experience on Arch Linux! 🎨✨** 