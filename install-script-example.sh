#!/bin/bash

# Install omarchy-monitor-settings for controlling Hyprland monitors
if [ -z "$OMARCHY_BARE" ] && ! command -v omarchy-monitor-settings &>/dev/null; then
  # Install Go if not already installed
  if ! command -v go &>/dev/null; then
    yay -S --noconfirm --needed go
  fi

  # Install the latest version using go install
  echo "Installing omarchy-monitor-settings..."
  go install github.com/ryanyogan/omarchy-monitor-settings@latest
  
  # Verify installation
  if command -v omarchy-monitor-settings &>/dev/null; then
    echo "✅ omarchy-monitor-settings installed successfully!"
    omarchy-monitor-settings --version
  else
    echo "❌ Installation failed. Trying alternative method..."
    
    # Fallback to building from source
    git clone https://github.com/ryanyogan/omarchy-monitor-settings.git /tmp/omarchy-monitor-settings
    cd /tmp/omarchy-monitor-settings
    
    # Checkout specific version for testing
    git checkout v1.1.2
    
    export CGO_ENABLED=0
    export GOOS=linux
    
    # Build with version information
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    go build -ldflags "-s -w -X main.version=${VERSION}" -o omarchy-monitor-settings ./cmd/omarchy-monitor-settings
    
    sudo mv omarchy-monitor-settings /usr/local/bin/
    sudo chmod +x /usr/local/bin/omarchy-monitor-settings
    cd -
    rm -rf /tmp/omarchy-monitor-settings
    
    echo "✅ omarchy-monitor-settings installed from source!"
  fi
fi 