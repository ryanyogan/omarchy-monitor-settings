# Maintainer: Your Name <your.email@example.com>
pkgname=hyprland-monitor-tui
pkgver=0.1.0
pkgrel=1
pkgdesc="A stunning TUI application for managing Hyprland monitor resolution and scaling"
arch=('x86_64' 'aarch64')
url="https://github.com/yourusername/hyprland-monitor-tui"
license=('MIT')
depends=()
optdepends=(
    'hyprland: Primary target window manager (recommended)'
    'wlr-randr: Wayland display manager fallback'
    'xorg-xrandr: X11 display manager fallback for compatibility'
)
makedepends=()
source=("$pkgname-$pkgver.tar.gz::file://$PWD")
sha256sums=('SKIP')

build() {
    cd "$srcdir"
    
    # Ensure Go is available via mise if not system-installed
    if ! command -v go &> /dev/null; then
        if command -v mise &> /dev/null; then
            eval "$(mise activate bash)"
            mise use go@latest
        else
            echo "Error: Go not found and mise not available"
            echo "Please install Go or mise before building"
            exit 1
        fi
    fi
    
    # Set Go environment
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64
    
    # Build the application
    go build -v \
        -buildmode=pie \
        -mod=readonly \
        -modcacherw \
        -ldflags "-linkmode external -extldflags \"${LDFLAGS}\" -s -w -X main.version=${pkgver}" \
        -o ${pkgname} .
}

package() {
    cd "$srcdir"
    
    # Install binary
    install -Dm755 "${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
    
    # Install desktop entry
    install -Dm644 "${pkgname}.desktop" "${pkgdir}/usr/share/applications/${pkgname}.desktop"
    
    # Install documentation
    install -Dm644 README.md "${pkgdir}/usr/share/doc/${pkgname}/README.md"
    
    # Install license
    install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
} 