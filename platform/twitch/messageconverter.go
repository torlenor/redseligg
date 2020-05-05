package twitch

import (
	"regexp"
	"strings"
)

func replaceAbyleBotterUserID(msg string) string {
	re := regexp.MustCompile(`(<@[a-zA-Z ]+>)`)

	foundUserIDs := map[string]bool{}

	matches := re.FindAllStringSubmatch(msg, -1)
	for _, match := range matches {
		if len(match) > 1 {
			foundUserIDs[match[1]] = true
		}
	}

	replaced := msg
	for userID := range foundUserIDs {
		replaced = strings.ReplaceAll(replaced, userID, userID[2:len(userID)-1])
	}

	return replaced
}

func convertMessageFromAbyleBotter(text string) string {
	return replaceAbyleBotterUserID(text)
}
