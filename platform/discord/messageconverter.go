package discord

import "strings"

func convertMessageFromAbyleBotter(text string) string {
	return strings.ReplaceAll(text, "\n", "\\n")
}
