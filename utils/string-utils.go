package utils

import (
	"regexp"
	"strings"
)

// IsNonEmptyString returns 'true' if string contains ANY
// non-whitespace characters
func IsNonEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) != 0
}

// IsEmptyString returns 'true' if string contains NO
// non-whitespace characters
func IsEmptyString(str string) bool {
	return !IsNonEmptyString(str)
}

// IsStringMissingInSlice returns 'true' if the provided slice
// DOES NOT contain provided string value
func IsStringMissingInSlice(str string, slice []string) bool {
	for _, element := range slice {
		if str == element {
			return false
		}
	}
	return true
}

// IsStringValidEmail returns 'true' if string meets requirements
// for an email address
func IsStringValidEmail(str string) bool {
	match, _ := regexp.MatchString("^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$", str)
	return match
}
