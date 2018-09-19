package echoplugin

import "strings"

func stripCmd(str string, cmd string) string {
	return strings.TrimPrefix(str, "!"+cmd+" ")
}
