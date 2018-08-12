package plugins

import "strings"

func stripCmd(str string, cmd string) string {
	return strings.TrimLeft(str, "!"+cmd)
}
