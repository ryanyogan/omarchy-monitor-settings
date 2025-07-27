package utils

import (
	"fmt"
	"strings"
)

// FormatScalePercent formats a scale as percentage.
func FormatScalePercent(scale float64) string {
	return fmt.Sprintf("%.0f%%", scale*100)
}

// FormatResolution formats width and height as resolution string.
func FormatResolution(width, height int) string {
	return fmt.Sprintf("%dx%d", width, height)
}

// FormatRefreshRate formats refresh rate with Hz suffix.
func FormatRefreshRate(rate float64) string {
	return fmt.Sprintf("%.2fHz", rate)
}

// TruncateString truncates string to max length with ellipsis.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// PadRight pads string to specified width.
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
