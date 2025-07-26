# Omarchy Monitor Settings

A terminal-based interface for managing monitor resolution and scaling in Hyprland/Wayland environments.

## Overview

This tool provides an interactive terminal interface for configuring monitor settings in Hyprland. It automatically detects connected displays and offers intelligent scaling recommendations based on resolution and DPI.

## Features

- Multi-monitor detection and configuration
- Intelligent scaling recommendations
- Real-time configuration updates
- Terminal-adaptive theming
- Demo mode for testing
- Comprehensive test coverage
- Git-based versioning system

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
go install github.com/ryanyogan/omarchy-monitor-settings@v1.1.1
```

### From Source

```bash
git clone https://github.com/ryanyogan/omarchy-monitor-settings.git
cd omarchy-monitor-settings
make build
sudo cp omarchy-monitor-settings /usr/local/bin/
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

# Development
make build        # Build binary
make clean        # Remove artifacts
make version      # Show version info
make help         # Show all targets
```

### Quality Assurance

The project maintains high code quality standards:

- Comprehensive test suite with 60+ tests
- Race condition detection
- Static analysis with `staticcheck`
- Code formatting validation
- Terminal theme adaptation testing

### Project Structure

```
├── main.go           # CLI entry point and configuration
├── model.go          # TUI model and rendering logic
├── monitor.go        # Monitor detection and configuration
├── scripts/          # Development scripts
│   └── version.sh    # Version management script
├── *_test.go         # Test files
├── Makefile          # Build automation
└── README.md         # Documentation
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

# Demo mode (no Hyprland required)
./omarchy-monitor-settings --no-hyprland-check
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

## Compatibility

- **Primary**: Arch Linux + Hyprland
- **Secondary**: Any Wayland compositor with `wlr-randr`
- **Development**: Any Unix-like system (demo mode)

## Contributing

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Run `make quality-check` before committing
4. Submit a pull request

### Code Standards

- All code must pass `make quality-check`
- Tests required for new functionality
- Follow Go best practices
- Maintain compatibility with Hyprland

## License

MIT License. See LICENSE file for details.

