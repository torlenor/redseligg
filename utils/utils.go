package utils

import (
	"fmt"
	"strings"
)

// StripCmd removes a bot command starting with ! from the string
func StripCmd(str string, cmd string) string {
	return strings.TrimPrefix(str, "!"+cmd+" ")
}

// GenerateErrorResponse generates a Json string with statusCode and statusMessage specified
func GenerateErrorResponse(statusCode uint16, statusText string) string {
	return fmt.Sprintf(`{"error": { "code": %d, "message": "%s" } }`, statusCode, statusText)
}
