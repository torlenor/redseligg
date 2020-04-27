package discord

import "strings"

func convertMessageFromAbyleBotter(text string) string {
	n := strings.ReplaceAll(text, "\n", "\\n")
	n = strings.ReplaceAll(n, `"`, `\"`)
	return n
}
