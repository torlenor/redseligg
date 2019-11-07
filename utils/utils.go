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
