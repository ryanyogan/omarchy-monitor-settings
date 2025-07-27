package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

// CommandExists checks if a command exists in PATH.
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// ParseResolution parses a resolution string like "1920x1080".
func ParseResolution(resStr string) (int, int) {
	parts := strings.Split(resStr, "x")
	if len(parts) != 2 {
		return 0, 0
	}

	width, err1 := strconv.Atoi(parts[0])
	height, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return 0, 0
	}

	return width, height
}

// ParseRefreshRate parses refresh rate from string like "60.00".
func ParseRefreshRate(refreshStr string) float64 {
	if rate, err := strconv.ParseFloat(refreshStr, 64); err == nil {
		return rate
	}
	return 60.0 // default fallback
}
