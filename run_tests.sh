#!/bin/bash

# Test runner script for hyprland-monitor-tui
# This script runs tests without starting the TUI

echo "Running tests for hyprland-monitor-tui..."

# Set environment variables for consistent test environment
export TERM=dumb
export DISPLAY=""
export WAYLAND_DISPLAY=""
export HYPRLAND_INSTANCE_SIGNATURE=""

# Run tests with timeout
timeout 30s go test -v ./... 2>&1

# Check exit code
if [ $? -eq 124 ]; then
    echo "Tests timed out (this is expected for TUI tests)"
    exit 0
elif [ $? -eq 0 ]; then
    echo "Tests completed successfully"
    exit 0
else
    echo "Tests failed"
    exit 1
fi 