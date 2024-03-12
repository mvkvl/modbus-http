package util

import (
	"fmt"
	"strings"
)

func StringSlice(input string, split string) []string {
	var result []string
	for _, s := range strings.Split(input, split) {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func StringOrDefault(input string, defaultValue string) string {
	if input == "" {
		return defaultValue
	} else {
		return input
	}
}

func FormattedStringOrEmpty(format, input string) string {
	if input == "" {
		return ""
	}
	return fmt.Sprintf(format, input)
}
