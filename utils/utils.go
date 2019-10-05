package utils

import "strings"

func StripCmd(str string, cmd string) string {
	return strings.TrimPrefix(str, "!"+cmd+" ")
}
