package utils

import (
	"fmt"
	"strings"
)

func FormatScalePercent(scale float64) string {
	return fmt.Sprintf("%.0f%%", scale*100)
}

func FormatResolution(width, height int) string {
	return fmt.Sprintf("%dx%d", width, height)
}

func FormatRefreshRate(rate float64) string {
	return fmt.Sprintf("%.2fHz", rate)
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
