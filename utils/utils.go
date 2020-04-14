package utils

import (
	"fmt"
	"strings"
)

// StripCmd removes a bot command starting with ! from the string
func StripCmd(str string, cmd string) string {
	return strings.TrimPrefix(str, "!"+cmd+" ")
}

// GenerateErrorResponse generates a Json error string
func GenerateErrorResponse(errorText string) string {
	return fmt.Sprintf(`{"error": "%s"}`, errorText)
}

// StringSliceContains check sif a string slice contains a certain entry and return true if yes.
// A nil slice will always return false.
func StringSliceContains(s []string, e string) bool {
	if s == nil {
		return false
	}

	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
