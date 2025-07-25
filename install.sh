#!/bin/bash

# Hyprland Monitor TUI Installation Script for Arch Linux
# This script installs the TUI and ensures all dependencies are available

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Pretty printing functions
print_header() {
    echo -e "${BLUE}================================================${NC}"
    echo -e "${CYAN}🖥️  Hyprland Monitor TUI Installation${NC}"
    echo -e "${BLUE}================================================${NC}"
}

print_step() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_success() {
    echo -e "${GREEN}🎉${NC} $1"
}

# Check if running on Arch Linux
check_arch() {
    if ! command -v pacman &> /dev/null; then
        print_error "This script is designed for Arch Linux systems with pacman."
        exit 1
    fi
    print_step "Detected Arch Linux system"
}

# Check if running on Hyprland
check_hyprland() {
    if [[ -n "$HYPRLAND_INSTANCE_SIGNATURE" ]]; then
        print_step "Running on Hyprland ✨"
        return 0
    fi
    
    if command -v hyprctl &> /dev/null; then
        print_info "Hyprland detected but not currently running"
        print_info "You can still install and test in demo mode"
        return 0
    fi
    
    print_info "Hyprland not detected in current session"
    print_info "This tool is designed for Hyprland but will install anyway"
    print_info "You can test the beautiful UI in demo mode"
}

# Install dependencies
install_dependencies() {
    print_step "Setting up build environment..."
    
    # Check if Go is available
    if ! command -v go &> /dev/null; then
        print_info "Go not found, installing via mise..."
        if ! command -v mise &> /dev/null; then
            print_error "mise not found. Please install mise first:"
            echo "  curl https://mise.run | sh"
            echo "  Or visit: https://mise.jdx.dev/getting-started.html"
            exit 1
        fi
        mise use -g go@latest
        print_step "Go installed via mise"
        
        # Reload shell environment to pick up Go
        export PATH="$HOME/.local/share/mise/installs/go/latest/bin:$PATH"
        
        # Verify Go is now available
        if ! command -v go &> /dev/null; then
            print_error "Go installation failed. Please check mise setup."
            exit 1
        fi
    else
        print_step "Go already available: $(go version)"
    fi
    
    # Optional dependencies for better monitor detection (don't fail if unavailable)
    print_info "Installing optional dependencies for better monitor detection..."
    sudo pacman -S --needed wlr-randr 2>/dev/null || print_info "wlr-randr not available, skipping"
    sudo pacman -S --needed xorg-xrandr 2>/dev/null || print_info "xrandr not available, skipping"
    
    print_step "Build environment ready"
}

# Build the application
build_app() {
    print_step "Building Hyprland Monitor TUI..."
    
    # Ensure we have the source
    if [[ ! -f "main.go" ]]; then
        print_error "Source files not found. Please run this script from the project directory."
        exit 1
    fi
    
    # Build with optimizations
    export CGO_ENABLED=0
    export GOOS=linux
    
    go build -v \
        -ldflags "-s -w -X main.version=1.0.0" \
        -o hyprland-monitor-tui .
    
    print_step "Build completed"
}

# Install the application
install_app() {
    print_step "Installing application..."
    
    # Install binary
    sudo install -Dm755 hyprland-monitor-tui /usr/local/bin/hyprland-monitor-tui
    
    # Install desktop entry if available
    if [[ -f "hyprland-monitor-tui.desktop" ]]; then
        sudo install -Dm644 hyprland-monitor-tui.desktop /usr/share/applications/hyprland-monitor-tui.desktop
        print_step "Desktop entry installed"
    fi
    
    # Install documentation
    if [[ -f "README.md" ]]; then
        sudo install -Dm644 README.md /usr/share/doc/hyprland-monitor-tui/README.md
        print_step "Documentation installed"
    fi
    
    print_step "Application installed to /usr/local/bin/hyprland-monitor-tui"
}

# Test the installation
test_installation() {
    print_step "Testing installation..."
    
    if command -v hyprland-monitor-tui &> /dev/null; then
        print_success "Installation successful!"
        echo
        print_info "You can now run the application with:"
        echo -e "  ${CYAN}hyprland-monitor-tui${NC}"
        echo
        if [[ -z "$HYPRLAND_INSTANCE_SIGNATURE" ]]; then
            print_info "Since you're not currently in Hyprland, you can test with:"
            echo -e "  ${CYAN}hyprland-monitor-tui --no-hyprland-check${NC}"
        fi
        echo
        print_info "The application features:"
        echo "  🎨 Beautiful Tokyo Night theme"
        echo "  🖥️  Smart monitor detection"
        echo "  📏 Intelligent scaling recommendations"
        echo "  ⚡ Real-time configuration"
        echo "  🎯 btop-like clean interface"
        echo
        print_info "Go version used: $(go version 2>/dev/null || echo 'managed by mise')"
        
    else
        print_error "Installation failed. Binary not found in PATH."
        exit 1
    fi
}

# Create uninstall script
create_uninstall() {
    print_step "Creating uninstall script..."
    
    cat > uninstall.sh << 'EOF'
#!/bin/bash
echo "🗑️  Uninstalling Hyprland Monitor TUI..."
sudo rm -f /usr/local/bin/hyprland-monitor-tui
sudo rm -f /usr/share/applications/hyprland-monitor-tui.desktop
sudo rm -rf /usr/share/doc/hyprland-monitor-tui/
echo "✓ Uninstallation complete"
EOF
    
    chmod +x uninstall.sh
    print_step "Uninstall script created (run ./uninstall.sh to remove)"
}

# Main installation process
main() {
    print_header
    
    check_arch
    check_hyprland
    install_dependencies
    build_app
    install_app
    test_installation
    create_uninstall
    
    echo
    print_success "🎉 Installation complete! Enjoy your beautiful TUI! 🎉"
    echo
}

# Run main function
main "$@" 