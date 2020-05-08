package discord

import "strings"

func convertMessageFromRedseligg(text string) string {
	n := strings.ReplaceAll(text, "\n", "\\n")
	n = strings.ReplaceAll(n, `"`, `\"`)
	return n
}
