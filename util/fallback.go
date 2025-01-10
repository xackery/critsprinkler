package util

import "strconv"

// FallbackInt returns the integer value of a string or a fallback value if the string is empty or cannot be converted to an integer.
func FallbackInt(in string, fallback int) int {
	if in == "" {
		return fallback
	}
	i, err := strconv.Atoi(in)
	if err != nil {
		return fallback
	}
	return i
}
