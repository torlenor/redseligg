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

// ExtractSubCommandAndArgsString returns the string splitted in the first word and the rest.
func ExtractSubCommandAndArgsString(message string) (subcommand string, argument string) {
	splitted := strings.Split(message, " ")
	if len(splitted) > 0 {
		subcommand = splitted[0]
	}
	if len(splitted) > 1 {
		argument = strings.Join(splitted[1:], " ")
	}
	return subcommand, argument
}
