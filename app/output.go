package app

import (
	"strings"
	"unicode"
)

func BytesToASCII(inp []byte) []byte {
	result := make([]byte, len(inp))
	for i, b := range inp {
		result[i] = byte(b % 127)
	}
	return result
}

func Coerce(inp []byte, length int) string {
	t := string(BytesToASCII(inp))
	t = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, t)
	t = strings.ReplaceAll(t, "\n", "")
	t = strings.ReplaceAll(t, " ", "")
	return t[:length]

}
