// Package types provides core data structures and constants for the monitor settings application.
package types

// Constants for UI layout - MUST MAINTAIN EXACT VALUES FOR UI CONSISTENCY
const (
	MinTerminalWidth  = 80
	MinTerminalHeight = 20
	HeaderHeight      = 7
	FooterHeight      = 5
)

// Scaling constants - MUST MAINTAIN EXACT VALUES FOR FUNCTIONALITY
const (
	MinMonitorScale = 0.5
	MaxMonitorScale = 4.0
	MinGTKScale     = 1
	MaxGTKScale     = 3
	MinFontDPI      = 72
	MaxFontDPI      = 300
	BaseDPI         = 96
)

// ValidHyprlandScales contains the valid scaling values for Hyprland.
var ValidHyprlandScales = []float64{
	1.0, 1.25, 1.33333, 1.5, 1.66667, 1.75, 2.0, 2.25, 2.5, 3.0,
}

// Error messages - MUST MAINTAIN EXACT TEXT FOR UI CONSISTENCY
const (
	ErrTerminalTooSmall = "Terminal too small\nPlease resize to at least 80x20"
	ErrNoMonitors       = "No monitors detected"
)
