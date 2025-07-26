#!/bin/bash

# Omarchy Monitor Settings Installation Script for Arch Linux
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
    echo -e "${CYAN}🖥️  Omarchy Monitor Settings Installation${NC}"
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
    # Removed xrandr support - using hyprctl and wlr-randr only
    
    print_step "Build environment ready"
}

# Build the application
build_app() {
    print_step "Building Omarchy Monitor Settings..."
    
    # Ensure we have the source
    if [[ ! -f "main.go" ]]; then
        print_error "Source files not found. Please run this script from the project directory."
        exit 1
    fi
    
    # Build with optimizations
    export CGO_ENABLED=0
    export GOOS=linux
    
    # Get version from git tag or use "dev"
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    
    go build -v \
        -ldflags "-s -w -X main.version=${VERSION}" \
        -o omarchy-monitor-settings .
    
    print_step "Build completed"
}

# Install the application
install_app() {
    print_step "Installing application..."
    
    # Install binary
    sudo install -Dm755 omarchy-monitor-settings /usr/local/bin/omarchy-monitor-settings
    
    # Install desktop entry if available
    if [[ -f "omarchy-monitor-settings.desktop" ]]; then
        sudo install -Dm644 omarchy-monitor-settings.desktop /usr/share/applications/omarchy-monitor-settings.desktop
        print_step "Desktop entry installed"
    fi
    
    # Install documentation
    if [[ -f "README.md" ]]; then
        sudo install -Dm644 README.md /usr/share/doc/omarchy-monitor-settings/README.md
        print_step "Documentation installed"
    fi
    
    print_step "Application installed to /usr/local/bin/omarchy-monitor-settings"
}

# Test the installation
test_installation() {
    print_step "Testing installation..."
    
    if command -v omarchy-monitor-settings &> /dev/null; then
        print_success "Installation successful!"
        echo
        print_info "You can now run the application with:"
        echo -e "  ${CYAN}omarchy-monitor-settings${NC}"
        echo
        if [[ -z "$HYPRLAND_INSTANCE_SIGNATURE" ]]; then
            print_info "Since you're not currently in Hyprland, you can test with:"
            echo -e "  ${CYAN}omarchy-monitor-settings --no-hyprland-check${NC}"
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
echo "🗑️  Uninstalling Omarchy Monitor Settings..."
sudo rm -f /usr/local/bin/omarchy-monitor-settings
sudo rm -f /usr/share/applications/omarchy-monitor-settings.desktop
sudo rm -rf /usr/share/doc/omarchy-monitor-settings/
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