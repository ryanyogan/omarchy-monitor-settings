# Omarchy Monitor Settings

Terminal UI for configuring monitor resolution and scaling on Hyprland/Wayland, with intelligent DPI recommendations.

## What It Is

Configuring monitors on Hyprland is a manual process involving `hyprctl` commands, fractional scaling math, and separate GTK/font DPI adjustments that interact in non-obvious ways. This tool wraps all of that into a Bubbletea TUI that detects connected displays, recommends scaling factors based on resolution and physical DPI, and applies changes across compositor, GTK, and font rendering in a single action.

It handles the edge cases that make monitor configuration annoying: Framework 13 laptops with unusual 2880x1800 panels, mixed-DPI multi-monitor setups, and the difference between compositor scaling and application-level scaling. Smart scaling presets (Ultra Sharp, Enhanced, Productive) encode tested recommendations for common resolutions from 1080p through 6K.

## Architecture

```
cmd/omarchy-monitor-settings/    # CLI entry (Cobra)
internal/
├── app/                          # Configuration management
├── monitor/                      # Detection via hyprctl/wlr-randr
└── tui/                          # Bubbletea model + rendering
pkg/
├── testing/                      # Visual regression framework
├── types/                        # Shared config types
├── ui/                           # Content + styles (Lipgloss)
└── utils/                        # Navigation, parsing, validation
```

Monitor detection cascades: `hyprctl` (Hyprland native) → `wlr-randr` (Wayland fallback) → demo data (development). Changes are applied immediately and reflected in the dashboard.

## Why This Matters

This is a supporting tool for the [Omarchy](https://github.com/ryanyogan/omarchy) Linux desktop environment. The problem it solves is specific but real: Wayland compositors handle fractional scaling differently than X11, and the interaction between compositor scale, GTK scale factor, and font DPI creates a combinatorial space that most users navigate by trial and error. The 80+ test suite includes visual regression tests against golden files to ensure the TUI renders correctly across terminal themes.

## Status

Stable. Published as an Arch Linux package. 16 GitHub stars. Used daily on Framework 13 and multi-monitor desktop setups.

## Stack

- **Language:** Go 1.19+
- **TUI:** Bubbletea + Lipgloss (Charm)
- **CLI:** Cobra
- **Terminal detection:** Termenv
- **Target:** Hyprland/Wayland (wlr-randr fallback)
- **Testing:** 80+ tests including visual regression
- **Packaging:** Arch Linux PKGBUILD, automated install script

## Installation

```bash
# Go install
go install github.com/ryanyogan/omarchy-monitor-settings@latest

# From source
git clone https://github.com/ryanyogan/omarchy-monitor-settings.git
cd omarchy-monitor-settings && make build
sudo cp omarchy-monitor-settings /usr/local/bin/

# Arch Linux
makepkg -si
```

## Usage

```bash
omarchy-monitor-settings                    # Normal operation
omarchy-monitor-settings --no-hyprland-check  # Demo mode (any system)
omarchy-monitor-settings --debug             # Debug output
```

Controls: `j/k` navigate, `Enter` select, `h/l` adjust manual values, `m` manual mode, `q` quit.

## License

MIT
