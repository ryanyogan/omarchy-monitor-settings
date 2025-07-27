// Package utils provides utility functions for mathematical calculations and validations.
package utils

import "math"

// CalculateEffectiveResolution calculates effective resolution after scaling.
func CalculateEffectiveResolution(width, height int, scale float64) (int, int) {
	effectiveWidth := int(float64(width) / scale)
	effectiveHeight := int(float64(height) / scale)
	return effectiveWidth, effectiveHeight
}

// CalculateScreenRealEstate calculates screen real estate percentage.
func CalculateScreenRealEstate(scale float64) float64 {
	return 100.0 / scale
}

// CalculateFontMultiplier calculates font scaling multiplier.
func CalculateFontMultiplier(dpi, baseDPI int) float64 {
	return float64(dpi) / float64(baseDPI)
}

// IsValidHyprlandScale checks if scale is a valid Hyprland scale.
func IsValidHyprlandScale(scale float64, validScales []float64) bool {
	for _, validScale := range validScales {
		if math.Abs(validScale-scale) < 0.001 {
			return true
		}
	}
	return false
}

// RoundToTwoDecimalPlaces rounds float to 2 decimal places.
func RoundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

// FindNextValidScale finds the next valid Hyprland scale.
func FindNextValidScale(current float64, up bool, validScales []float64) float64 {
	currentIndex := -1
	for i, scale := range validScales {
		if math.Abs(scale-current) < 0.001 {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		minDiff := math.Abs(validScales[0] - current)
		currentIndex = 0
		for i, scale := range validScales {
			diff := math.Abs(scale - current)
			if diff < minDiff {
				minDiff = diff
				currentIndex = i
			}
		}
	}

	if up {
		if currentIndex < len(validScales)-1 {
			return validScales[currentIndex+1]
		}
		return validScales[len(validScales)-1]
	} else {
		if currentIndex > 0 {
			return validScales[currentIndex-1]
		}
		return validScales[0]
	}
}
