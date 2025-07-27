# 🏗️ Arch Linux Package Installation Guide

## 🚀 Quick Installation (Recommended)

### Method 1: Automated Installation Script

```bash
# Download and run the installation script
curl -sSL https://github.com/yourusername/omarchy-monitor-settings/raw/main/install.sh | bash

# Or if you have the source:
./install.sh
```

### Method 2: Manual Build & Install

```bash
# Install Go via mise (if not already available)
mise use -g go@latest

# Install optional dependencies (recommended)
sudo pacman -S wlr-randr

# Build the application
go build -o omarchy-monitor-settings .

# Install system-wide
sudo install -Dm755 omarchy-monitor-settings /usr/local/bin/omarchy-monitor-settings
sudo install -Dm644 omarchy-monitor-settings.desktop /usr/share/applications/omarchy-monitor-settings.desktop
```

## 📦 Package Creation for AUR

### Create Package Archive

```bash
# Create source archive
tar -czf omarchy-monitor-settings-1.1.0.tar.gz \
    main.go model.go monitor.go go.mod go.sum \
    README.md LICENSE omarchy-monitor-settings.desktop

# Generate checksums
sha256sum omarchy-monitor-settings-1.1.0.tar.gz
```

### Build with makepkg

```bash
# Build package
makepkg -si

# Or just build without installing
makepkg

# Install built package
sudo pacman -U omarchy-monitor-settings-1.1.0-1-x86_64.pkg.tar.zst
```

## 🔧 Arch Linux Compatibility

### ✅ Verified Components

- **Hyprland Integration**: Native `hyprctl` support
- **Wayland Support**: Uses `wlr-randr` for fallback detection
# Removed X11/xrandr support
- **Go Runtime**: Built with Go 1.21+ (available in Arch repos)
- **Dependencies**: All optional deps available in official repos

### 📋 System Requirements

```bash
# Build requirements (managed via mise)
mise use -g go@latest

# Runtime requirements (assumed to be installed)
pacman -S hyprland

# Optional but recommended
pacman -S wlr-randr

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
omarchy-monitor-settings

# Debug mode for troubleshooting
omarchy-monitor-settings --debug
```

### Outside Hyprland (Testing)

```bash
# Demo mode with sample data
omarchy-monitor-settings --no-hyprland-check

# This will show realistic Framework 13 + external monitor data
# Perfect for UI testing and development
```

### Verification Commands

```bash
# Check monitor detection methods
hyprctl monitors  # Should work in Hyprland
wlr-randr         # Should work in any Wayland compositor
# xrandr removed - no longer supported

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
   sudo pacman -S wlr-randr
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
which hyprctl wlr-randr

# Test monitor detection
omarchy-monitor-settings --debug --no-hyprland-check
```

## 🏷️ Package Information

- **Package Name**: `omarchy-monitor-settings`
- **Version**: `1.1.0`
- **Architecture**: `x86_64`, `aarch64`
- **License**: `MIT`
- **Dependencies**: `hyprland`
- **Optional Dependencies**: `wlr-randr`
- **Build Dependencies**: `go`

## 🗑️ Uninstallation

```bash
# If installed via script
./uninstall.sh

# If installed via pacman
sudo pacman -R omarchy-monitor-settings

# Manual removal
sudo rm /usr/local/bin/omarchy-monitor-settings
sudo rm /usr/share/applications/omarchy-monitor-settings.desktop
sudo rm -rf /usr/share/doc/omarchy-monitor-settings/
```

---

**Ready to create the most beautiful monitor management experience on Arch Linux! 🎨✨** 