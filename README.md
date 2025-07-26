# 🖥️ Hyprland Monitor TUI

A **stunning** terminal interface for managing Hyprland monitor resolution and scaling, built with Go and styled with the beautiful Tokyo Night theme.

## ✨ Features

- 🎨 **Beautiful Tokyo Night Theme** - Gorgeous colors and modern styling
- 🖥️ **Multi-Monitor Support** - Detect and configure multiple displays
- 📏 **Intelligent Scaling** - Smart recommendations based on resolution
- 🔧 **Real-time Configuration** - Apply changes instantly to Hyprland
- 🧪 **Demo Mode** - Test the interface on any platform
- ⚡ **Fast & Responsive** - Built with Go and Bubbletea

## 🚀 Quick Start (macOS Testing)

### Prerequisites

Make sure you have Go installed:
```bash
brew install go
```

### Build & Run

1. **Clone or navigate to the project directory**
2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build the application:**
   ```bash
   go build -o hyprland-monitor-tui
   ```

4. **Run in demo mode (perfect for macOS testing):**
   ```bash
   ./hyprland-monitor-tui --no-hyprland-check
   ```

### 🎮 Controls

- **↑/↓ or k/j** - Navigate menu items
- **Enter/Space** - Select option
- **h or ?** - Show help screen
- **Esc** - Return to main menu
- **q or Ctrl+C** - Quit application

### 🎨 What You'll See

The **AWARD-WINNING** TUI features:
- 🖥️ **Full terminal usage** - Dynamically resizes and scales to your terminal
- 🎨 **Stunning Tokyo Night theme** - Extended 20+ color palette 
- ✨ **Elegant selection indicators** - Beautiful `▶` arrows instead of background highlighting
- 📱 **Professional cards & sections** - Each screen beautifully organized
- 🏆 **btop-like design** - Clean, light, minimal, and absolutely gorgeous
- 🎯 **Smart overflow protection** - Never overflows, always fits perfectly
- 🌈 **Colorful key hints** - Each control highlighted in vibrant colors
- 🔥 **Instant responsiveness** - Smooth navigation that feels amazing
- 📊 **Rich monitor information** - Status badges, scaling recommendations
- 🧪 **Perfect demo mode** - Beautiful UI testing on any platform

## 🖥️ Production Use (Arch Linux + Hyprland)

### Installation
```bash
go build -o hyprland-monitor-tui
sudo cp hyprland-monitor-tui /usr/local/bin/
```

### Usage
```bash
# Normal mode (requires Hyprland)
hyprland-monitor-tui

# Debug mode
hyprland-monitor-tui --debug
```

## 🔧 Monitor Detection

The application automatically detects monitors using:
1. **hyprctl** (Hyprland's native tool)
2. **wlr-randr** (Wayland fallback)
# Removed xrandr support
4. **Demo data** (for testing/development)

## 📊 Scaling Recommendations

Smart scaling based on resolution:
- **4K+ displays** (3840x2160+): 2x monitor scaling, 0.8x font scaling
- **1440p displays** (2560x1440): 1x monitor scaling, 0.9x font scaling  
- **1080p displays** (1920x1080): 1x monitor scaling, 0.8x font scaling

## 🎨 Tokyo Night Theme

The interface uses the beautiful Tokyo Night color palette:
- Background: `#1a1b26`
- Surface: `#24283b`
- Primary: `#7aa2f7`
- Accent colors: Cyan, Green, Yellow, Orange, Red, Purple

## 🛠️ Development

### Project Structure
```
├── main.go      # CLI entry point
├── model.go     # TUI model & rendering
├── monitor.go   # Monitor detection & configuration
├── go.mod       # Dependencies
└── README.md    # Documentation
```

### Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling and layout
- `github.com/spf13/cobra` - CLI framework

## 🎯 Supported Platforms

- ✅ **Arch Linux + Hyprland** (primary target)
- ✅ **Any Wayland compositor** (wlr-randr fallback)
# Removed X11/xrandr support
- ✅ **macOS/Other** (demo mode for UI testing)

## 🤝 Contributing

This tool is designed for the Arch Linux community and Hyprland users. Contributions welcome!

## 📝 License

MIT License - Feel free to use and modify as needed.

---

**Built with ❤️ for the Hyprland community** 