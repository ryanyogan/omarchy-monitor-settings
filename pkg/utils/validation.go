package utils

func ValidateGTKScale(scale, minVal, maxVal int) int {
	if scale < minVal {
		return minVal
	}
	if scale > maxVal {
		return maxVal
	}
	return scale
}

func ValidateMonitorScale(scale, minVal, maxVal float64) float64 {
	if scale < minVal {
		return minVal
	}
	if scale > maxVal {
		return maxVal
	}
	return scale
}

func ValidateFontDPI(dpi, minVal, maxVal int) int {
	if dpi < minVal {
		return minVal
	}
	if dpi > maxVal {
		return maxVal
	}
	return dpi
}
