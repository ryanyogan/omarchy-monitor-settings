package utils

// Navigation utilities - PURE LOGIC, NO UI IMPACT

// ClampIndex ensures an index stays within bounds
func ClampIndex(index, min, max int) int {
	if index < min {
		return min
	}
	if index > max {
		return max
	}
	return index
}

// WrapIndex wraps an index around bounds (for circular navigation)
func WrapIndex(index, min, max int) int {
	if max < min {
		return min
	}

	if index < min {
		return max
	}
	if index > max {
		return min
	}
	return index
}

// NavigateUp decreases index with bounds checking
func NavigateUp(currentIndex, minIndex int) int {
	if currentIndex > minIndex {
		return currentIndex - 1
	}
	return currentIndex
}

// NavigateDown increases index with bounds checking
func NavigateDown(currentIndex, maxIndex int) int {
	if currentIndex < maxIndex {
		return currentIndex + 1
	}
	return currentIndex
}

// IsValidIndex checks if an index is within bounds
func IsValidIndex(index, min, max int) bool {
	return index >= min && index <= max
}

// GetSafeIndex returns a safe index within bounds, defaulting to min if invalid
func GetSafeIndex(index, min, max int) int {
	if IsValidIndex(index, min, max) {
		return index
	}
	return min
}
