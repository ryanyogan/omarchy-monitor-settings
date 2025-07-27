package utils

// ValidateGTKScale ensures GTK scale is within bounds.
func ValidateGTKScale(scale, min, max int) int {
	if scale < min {
		return min
	}
	if scale > max {
		return max
	}
	return scale
}

// ValidateMonitorScale ensures monitor scale is valid.
func ValidateMonitorScale(scale, min, max float64) float64 {
	if scale < min {
		return min
	}
	if scale > max {
		return max
	}
	return scale
}

// ValidateFontDPI ensures font DPI is within bounds.
func ValidateFontDPI(dpi, min, max int) int {
	if dpi < min {
		return min
	}
	if dpi > max {
		return max
	}
	return dpi
}
