package utils

import (
	"strconv"
	"strings"
)

func ContainsAll(text string, substrings ...string) bool {
	for _, substr := range substrings {
		if !strings.Contains(text, substr) {
			return false
		}
	}
	return true
}

func ContainsAny(text string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(text, substr) {
			return true
		}
	}
	return false
}

func ExtractFloat64(s string, defaultValue float64) float64 {
	if value, err := strconv.ParseFloat(s, 64); err == nil {
		return value
	}
	return defaultValue
}

func ExtractInt(s string, defaultValue int) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	}
	return defaultValue
}

func SplitAndTrim(text, separator string) []string {
	parts := strings.Split(text, separator)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

func IsNonEmptyLine(line string) bool {
	return strings.TrimSpace(line) != ""
}

func HasPrefixIgnoreCase(text, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(text), strings.ToLower(prefix))
}

func SafeStringAccess(slice []string, index int) string {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	return ""
}
