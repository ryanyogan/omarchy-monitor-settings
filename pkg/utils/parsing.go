package utils

import (
	"strconv"
	"strings"
)

// Parsing utilities - PURE LOGIC, NO UI IMPACT

// ContainsAll checks if a string contains all specified substrings
func ContainsAll(text string, substrings ...string) bool {
	for _, substr := range substrings {
		if !strings.Contains(text, substr) {
			return false
		}
	}
	return true
}

// ContainsAny checks if a string contains any of the specified substrings
func ContainsAny(text string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(text, substr) {
			return true
		}
	}
	return false
}

// ExtractFloat64 extracts a float64 from a string, returns default value if failed
func ExtractFloat64(s string, defaultValue float64) float64 {
	if value, err := strconv.ParseFloat(s, 64); err == nil {
		return value
	}
	return defaultValue
}

// ExtractInt extracts an int from a string, returns default value if failed
func ExtractInt(s string, defaultValue int) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	}
	return defaultValue
}

// SplitAndTrim splits a string and trims whitespace from each part
func SplitAndTrim(text, separator string) []string {
	parts := strings.Split(text, separator)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// IsNonEmptyLine checks if a line is not empty or just whitespace
func IsNonEmptyLine(line string) bool {
	return strings.TrimSpace(line) != ""
}

// HasPrefixIgnoreCase checks if string has prefix ignoring case
func HasPrefixIgnoreCase(text, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(text), strings.ToLower(prefix))
}

// SafeStringAccess safely accesses a string slice at index, returns empty string if out of bounds
func SafeStringAccess(slice []string, index int) string {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	return ""
}
