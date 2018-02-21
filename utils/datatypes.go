package utils

import (
	"strconv"
	"strings"
)

// strIsInt checks if string can be parsed as an integer
func StrIsInt(text string) bool {
	// Attempt to parse string as int
	if _, err := strconv.Atoi(text); err == nil {
		return true
	}
	return false
}

// strIsFloat checks if string cam be parsed as a float
func StrIsFloat(text string) bool {
	// Attempt to parse string as float64
	if _, err := strconv.ParseFloat(text, 64); err == nil {
		// Check for "." character
		return strings.Contains(text, ".")
	}
	return false
}
