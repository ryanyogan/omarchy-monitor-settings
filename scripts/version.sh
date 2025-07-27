#!/bin/bash

# Version management script for omarchy-monitor-settings

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Get current version from git
get_current_version() {
    git describe --tags --always --dirty 2>/dev/null || echo "dev"
}

# Get latest tag
get_latest_tag() {
    git describe --tags --abbrev=0 2>/dev/null || echo "none"
}

# Get current commit
get_current_commit() {
    git rev-parse --short HEAD 2>/dev/null || echo "unknown"
}

# Check if working directory is clean
is_clean() {
    git diff-index --quiet HEAD -- 2>/dev/null
}

# Show version information
show_version() {
    echo "Version Information:"
    echo "==================="
    echo "Current version: $(get_current_version)"
    echo "Latest tag: $(get_latest_tag)"
    echo "Current commit: $(get_current_commit)"
    echo "Working directory clean: $(is_clean && echo "yes" || echo "no")"
}

# Create a new version tag
create_tag() {
    local version=$1
    
    if [[ -z "$version" ]]; then
        print_warning "Please provide a version number (e.g., 1.2.0)"
        exit 1
    fi
    
    if ! is_clean; then
        print_warning "Working directory is not clean. Please commit or stash changes first."
        exit 1
    fi
    
    print_info "Creating tag v${version}..."
    git tag -a "v${version}" -m "Release v${version}"
    print_success "Tag v${version} created successfully!"
    
    print_info "To push the tag to remote:"
    echo "  git push origin v${version}"
}

# Build with current version
build_with_version() {
    local version=$(get_current_version)
    print_info "Building with version: ${version}"
    
    go build -ldflags "-X main.version=${version}" -o omarchy-monitor-settings ./cmd/omarchy-monitor-settings
    print_success "Build completed successfully!"
}

# Main script logic
case "${1:-show}" in
    "show")
        show_version
        ;;
    "tag")
        create_tag "$2"
        ;;
    "build")
        build_with_version
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command] [version]"
        echo ""
        echo "Commands:"
        echo "  show    - Show current version information (default)"
        echo "  tag     - Create a new version tag (requires version number)"
        echo "  build   - Build the application with current version"
        echo "  help    - Show this help message"
        echo ""
        echo "Examples:"
        echo "  $0                    # Show version info"
        echo "  $0 tag 1.2.0         # Create tag v1.2.0"
        echo "  $0 build             # Build with current version"
        ;;
    *)
        print_warning "Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac 