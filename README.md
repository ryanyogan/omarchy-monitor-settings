# Omarchy Monitor Settings

A terminal-based interface for managing monitor resolution and scaling in Hyprland/Wayland environments.

## Overview

This tool provides an interactive terminal interface for configuring monitor settings in Hyprland. It automatically detects connected displays and offers intelligent scaling recommendations based on resolution and DPI.

## Features

- **Multi-monitor detection and configuration**
- **Intelligent scaling recommendations** with resolution-specific optimizations
- **High DPI display support** for 2.8K, 4K, 5K, and 6K displays
- **Framework 13 optimized** with specialized 2.8K scaling options
- **Real-time configuration updates** with dashboard reflection
- **Terminal-adaptive theming**
- **Demo mode for testing**
- **Comprehensive test coverage** (80+ tests)
- **Visual regression testing** for UI consistency
- **Git-based versioning system**
- **Automated installation and deployment**
- **Arch Linux package support**
- **Modular Go architecture**

## Requirements

- Go 1.19 or later
- For production use: Hyprland with `hyprctl` or `wlr-randr`
- For development: Any Unix-like system

## Installation

### Quick Install (Recommended)

```bash
# Install the latest version
go install github.com/ryanyogan/omarchy-monitor-settings@latest

# Install a specific version
go install github.com/ryanyogan/omarchy-monitor-settings@v1.1.2

# The application can now be installed directly from the root module
```

### From Source

```bash
git clone https://github.com/ryanyogan/omarchy-monitor-settings.git
cd omarchy-monitor-settings
make build
sudo cp omarchy-monitor-settings /usr/local/bin/
```

### Arch Linux Package

```bash
# Using the provided PKGBUILD
makepkg -si

# Or install from AUR (when available)
yay -S omarchy-monitor-settings
```

### Automated Installation Script

```bash
# Run the automated installation script
./install.sh

# This script will:
# - Install Go if needed (via mise)
# - Install optional dependencies (wlr-randr)
# - Build the application
# - Install to /usr/local/bin/
# - Create desktop entry
# - Create uninstall script
```

### Dependencies

The application will use available monitor detection tools in order of preference:
1. `hyprctl` (Hyprland native)
2. `wlr-randr` (Wayland fallback)
3. Demo data (development/testing)

## Usage

### Basic Usage

```bash
# Normal operation (requires Hyprland/Wayland)
omarchy-monitor-settings

# Demo mode (works on any system)
omarchy-monitor-settings --no-hyprland-check

# Debug mode
omarchy-monitor-settings --debug
```

### Controls

- `↑/↓` or `k/j` - Navigate menus
- `←/→` or `h/l` - Adjust values (manual scaling)
- `Enter/Space` - Select option
- `m` - Switch to manual scaling
- `h` or `?` - Help screen
- `Esc` - Return to previous screen
- `q` or `Ctrl+C` - Quit

### Smart Scaling Options

The application provides intelligent scaling recommendations based on your monitor's resolution and DPI:

#### High DPI Displays (2.8K, 4K, 5K, 6K)
- **2x Ultra Sharp**: Perfect integer scaling for maximum clarity
- **1.75x Enhanced**: Great balance of clarity and screen space
- **1.5x Productive**: Maximum screen real estate for workflows

#### Framework 13 (2.8K Display)
- **2x Ultra Sharp**: Recommended for Framework 13 with 2880x1800 display
- **1.75x Enhanced**: Excellent productivity balance
- **1.5x Productive**: Maximum workspace for development

#### Standard Displays (1080p, 2.5K)
- **1x Native**: Standard scaling for most use cases
- **1.25x Enhanced**: Slightly larger text for better readability
- **1.5x Large**: Accessibility-friendly larger text

## Versioning

The application uses Git-based versioning with build-time variable injection. This is the idiomatic Go approach for version management.

### Version Information

```bash
# Show current version info
make version

# Or use the version script
./scripts/version.sh
```

### Version Sources

- **Tagged releases**: Uses Git tag (e.g., `v1.1.0`)
- **Development builds**: Uses commit hash with `-dirty` suffix if uncommitted changes
- **Fallback**: Uses `dev` if Git is not available

### Creating Releases

```bash
# Create a new version tag
./scripts/version.sh tag 1.2.0

# Build with current version
./scripts/version.sh build

# Or use make
make build
```

### Version Output Examples

```bash
# Tagged release
omarchy-monitor-settings version v1.1.0

# Development build
omarchy-monitor-settings version c34ff8c-dirty

# Clean development build
omarchy-monitor-settings version c34ff8c
```

## Development

### Build System

The project uses a Makefile for development tasks:

```bash
# Run comprehensive quality checks
make quality-check

# Individual checks
make vet          # Code analysis
make fmt-check    # Format validation
make test-race    # Race condition detection
make staticcheck  # Static analysis
make build-check  # Compilation verification

# Testing
make test         # Basic test suite
make test-verbose # Detailed test output
make test-coverage # Coverage report

# Visual Regression Testing
make visual-test  # Run visual regression tests
make visual-update # Update golden files
make visual-clean # Clean test artifacts

# Development
make build        # Build binary
make clean        # Remove artifacts
make version      # Show version info
make help         # Show all targets
```

### Quality Assurance

The project maintains high code quality standards:

- Comprehensive test suite with 80+ tests
- Race condition detection
- Static analysis with `staticcheck`
- Code formatting validation
- Terminal theme adaptation testing
- Visual regression testing for UI consistency
- Automated installation and deployment scripts

### Project Structure

```
├── cmd/
│   └── omarchy-monitor-settings/  # Main application entry point
│       ├── main.go                # CLI entry point and configuration
│       └── main_test.go           # Main application tests
├── internal/
│   ├── app/                       # Application services and configuration
│   │   └── config.go              # Configuration management
│   ├── monitor/                   # Monitor detection and management
│   │   ├── monitor.go             # Monitor detection and configuration
│   │   └── monitor_test.go        # Monitor tests
│   └── tui/                       # Terminal user interface
│       ├── model.go               # TUI model and rendering logic
│       ├── model_test.go          # TUI unit tests
│       └── visual_regression_test.go # Visual regression tests
├── pkg/                           # Public packages
│   ├── testing/                   # Testing utilities
│   │   └── visual.go              # Visual testing framework
│   ├── types/                     # Shared types and constants
│   │   └── config.go              # Configuration types
│   ├── ui/                        # UI components
│   │   ├── content.go             # Content rendering
│   │   └── styles.go              # Styling definitions
│   └── utils/                     # Utility functions
│       ├── math.go                # Mathematical utilities
│       ├── navigation.go          # Navigation helpers
│       ├── parsing.go             # Text parsing utilities
│       ├── system.go              # System interaction
│       ├── text.go                # Text formatting
│       └── validation.go          # Validation functions
├── scripts/                       # Development and deployment scripts
│   └── version.sh                 # Version management script
├── testdata/                      # Test data and golden files
│   └── golden/                    # Visual regression test baselines
├── docs/                          # Documentation
│   ├── arch-package.md            # Arch Linux packaging guide
│   ├── hyprland-compatibility.md  # Hyprland compatibility notes
│   └── VISUAL_TESTING.md          # Visual testing documentation
├── install.sh                     # Automated installation script
├── install-script-example.sh      # Example installation script
├── debug-monitors.sh              # Monitor debugging script
├── PKGBUILD                       # Arch Linux package definition
├── omarchy-monitor-settings.desktop # Desktop entry file
├── Makefile                       # Build automation
└── README.md                      # This documentation
```

### Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/spf13/cobra` - CLI framework
- `github.com/muesli/termenv` - Terminal detection

## Testing

### Local Testing

```bash
# Run all tests
make test-verbose

# Test with race detection
make test-race

# Visual regression testing
make visual-test

# Demo mode (no Hyprland required)
./omarchy-monitor-settings --no-hyprland-check

# Debug monitor detection
./debug-monitors.sh
```

### CI/CD Integration

```bash
# Complete quality check (suitable for CI)
make quality-check
```

## Configuration

### Scaling Options

The application provides intelligent scaling recommendations:

- **4K+ displays** (3840x2160+): Higher scaling factors
- **1440p displays** (2560x1440): Moderate scaling
- **1080p displays** (1920x1080): Minimal scaling

### Manual Configuration

Users can manually adjust:
- Monitor scale (compositor-level)
- GTK scale (application-level)
- Font DPI (text rendering)

### Real-time Updates

Changes are applied immediately and reflected in the dashboard:
- Smart scaling: Applies recommended settings
- Manual scaling: Applies custom values
- Dashboard updates show current scale values

## Compatibility

- **Primary**: Arch Linux + Hyprland
- **Secondary**: Any Wayland compositor with `wlr-randr`
- **Development**: Any Unix-like system (demo mode)

## Contributing

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Run `make quality-check` before committing
4. Run visual regression tests: `make visual-test`
5. Submit a pull request

### Code Standards

- All code must pass `make quality-check`
- Tests required for new functionality
- Visual regression tests for UI changes
- Follow Go best practices
- Maintain compatibility with Hyprland
- Use the new modular structure (`cmd/`, `internal/`, `pkg/`)

### Architecture

The project follows Go best practices with a modular structure:
- `cmd/`: Application entry points
- `internal/`: Private application code
- `pkg/`: Public packages for reuse
- `scripts/`: Development and deployment utilities

## Deployment

### Automated Deployment

The project includes several deployment options:

#### Arch Linux Package
```bash
# Build and install package
makepkg -si

# Package includes:
# - Binary in /usr/bin/
# - Desktop entry in /usr/share/applications/
# - Documentation in /usr/share/doc/
# - License in /usr/share/licenses/
```

#### Automated Installation Script
```bash
# Full automated installation
./install.sh

# Features:
# - Go installation via mise (if needed)
# - Dependency installation (wlr-randr)
# - Application building and installation
# - Desktop entry creation
# - Uninstall script generation
```

#### Manual Installation
```bash
# Build and install manually
make build
sudo install -Dm755 omarchy-monitor-settings /usr/local/bin/
sudo install -Dm644 omarchy-monitor-settings.desktop /usr/share/applications/
```

### Development Scripts

#### Version Management
```bash
# Show version information
./scripts/version.sh

# Create new version tag
./scripts/version.sh tag 1.2.0

# Build with current version
./scripts/version.sh build
```

#### Debugging
```bash
# Debug monitor detection
./debug-monitors.sh

# This script will:
# - Check environment variables
# - Verify available commands
# - Test hyprctl output
# - Run the application in debug mode
```

#### Example Installation
```bash
# Example installation script for automation
./install-script-example.sh

# Features:
# - Go installation via package manager
# - go install fallback
# - Source build fallback
# - Verification and testing
```

## License

MIT License. See LICENSE file for details.

