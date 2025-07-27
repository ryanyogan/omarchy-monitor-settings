package utils

func ClampIndex(index, minVal, maxVal int) int {
	if index < minVal {
		return minVal
	}
	if index > maxVal {
		return maxVal
	}
	return index
}

func WrapIndex(index, minVal, maxVal int) int {
	if maxVal < minVal {
		return minVal
	}

	if index < minVal {
		return maxVal
	}
	if index > maxVal {
		return minVal
	}
	return index
}

func NavigateUp(currentIndex, minIndex int) int {
	if currentIndex > minIndex {
		return currentIndex - 1
	}
	return currentIndex
}

func NavigateDown(currentIndex, maxIndex int) int {
	if currentIndex < maxIndex {
		return currentIndex + 1
	}
	return currentIndex
}

func IsValidIndex(index, minVal, maxVal int) bool {
	return index >= minVal && index <= maxVal
}

func GetSafeIndex(index, minVal, maxVal int) int {
	if IsValidIndex(index, minVal, maxVal) {
		return index
	}
	return minVal
}
