=== Hyprland Monitor Detection Debug ===

1. Environment Check:
   HYPRLAND_INSTANCE_SIGNATURE: 4e242d086e20b32951fdc0ebcbfb4d41b5be8dcc_1753301308_799328318
   XDG_CURRENT_DESKTOP: Hyprland
   WAYLAND_DISPLAY: wayland-1

2. Available Commands:
   ✓ hyprctl: /usr/bin/hyprctl
   ✗ wlr-randr: not found
   # xrandr support removed

3. hyprctl monitors output:
   --- Raw Output ---
   Monitor eDP-1 (ID 0):
   2880x1920@120.00000 at 0x0
   description: BOE NE135A1M-NY1
   make: BOE
   model: NE135A1M-NY1
   serial:
   active workspace: 5 (5)
   special workspace: 0 ()
   reserved: 0 26 0 0
   scale: 1.67
   transform: 0
   focused: yes
   dpmsStatus: 1
   vrr: false
   solitary: 0
   activelyTearing: false
   directScanoutTo: 0
   disabled: false
   currentFormat: XRGB8888
   mirrorOf: none
   availableModes: 2880x1920@120.00Hz 2880x1920@60.00Hz 1920x1200@120.00Hz 1920x1080@120.00Hz 1600x1200@120.00Hz 1680x1050@120.00Hz 1280x1024@120.00Hz 1440x900@120.00Hz 1280x800@120.00Hz 1280x720@120.00Hz 1024x768@120.00Hz 800x600@120.00Hz 640x480@120.00Hz

--- End Output ---

Exit code: 0

4. Running TUI with debug:
   Running: ./hyprland-monitor-tui --debug
   ================================================
   DEBUG: Starting monitor detection...
   DEBUG: Trying method 1: hyprctl
   DEBUG: Found hyprctl, running 'hyprctl monitors'
   DEBUG: hyprctl output (634 bytes):
   Monitor eDP-1 (ID 0):
   2880x1920@120.00000 at 0x0
   description: BOE NE135A1M-NY1
   make: BOE
   model: NE135A1M-NY1
   serial:
   active workspace: 5 (5)
   special workspace: 0 ()
   reserved: 0 26 0 0
   scale: 1.67
   transform: 0
   focused: yes
   dpmsStatus: 1
   vrr: false
   solitary: 0
   activelyTearing: false
   directScanoutTo: 0
   disabled: false
   currentFormat: XRGB8888
   mirrorOf: none
   availableModes: 2880x1920@120.00Hz 2880x1920@60.00Hz 1920x1200@120.00Hz 1920x1080@120.00Hz 1600x1200@120.00Hz 1680x1050@120.00Hz 1280x1024@120.00Hz 1440x900@120.00Hz 1280x800@120.00Hz 1280x720@120.00Hz 1024x768@120.00Hz 800x600@120.00Hz 640x480@120.00Hz

DEBUG: Parsed 1 monitors from hyprctl output, parse error: <nil>
DEBUG: Monitor 0: eDP-1 (0x0@0.0Hz, scale 1.7)
DEBUG: Method returned 1 monitors, error: <nil>
DEBUG: Successfully detected 1 monitors using hyprctl
DEBUG: DetectMonitors returned 1 monitors, error: <nil>
DEBUG: Setting live mode - detected real monitors
DEBUG: Final state - isDemoMode: false, monitor count: 1
