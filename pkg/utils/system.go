package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

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

func ParseRefreshRate(refreshStr string) float64 {
	if rate, err := strconv.ParseFloat(refreshStr, 64); err == nil {
		return rate
	}
	return 60.0
}
