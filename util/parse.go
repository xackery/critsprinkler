package util

import "regexp"

// Parse parses a line with a regex and size
func Parse(line string, regex string, size int) ([]string, bool) {
	match := regexp.MustCompile(regex).FindStringSubmatch(line)
	if len(match) < 1 {
		return nil, false
	}
	match = match[1:]
	if len(match) != size {
		return nil, false
	}

	return match, true
}
