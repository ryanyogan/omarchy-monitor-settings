#!/bin/bash
echo "=== Hyprland Monitor Detection Debug ==="
echo

echo "1. Environment Check:"
echo "HYPRLAND_INSTANCE_SIGNATURE: ${HYPRLAND_INSTANCE_SIGNATURE:-'(not set)'}"
echo "XDG_CURRENT_DESKTOP: ${XDG_CURRENT_DESKTOP:-'(not set)'}"
echo "WAYLAND_DISPLAY: ${WAYLAND_DISPLAY:-'(not set)'}"
echo

echo "2. Available Commands:"
for cmd in hyprctl wlr-randr xrandr; do
    if command -v "$cmd" &> /dev/null; then
        echo "✓ $cmd: $(which $cmd)"
    else
        echo "✗ $cmd: not found"
    fi
done
echo

echo "3. hyprctl monitors output:"
if command -v hyprctl &> /dev/null; then
    echo "--- Raw Output ---"
    hyprctl monitors 2>&1
    echo "--- End Output ---"
    echo
    echo "Exit code: $?"
else
    echo "hyprctl not available"
fi
echo

echo "4. Running TUI with debug:"
echo "Running: ./hyprland-monitor-tui --debug"
echo "================================================"
./hyprland-monitor-tui --debug 