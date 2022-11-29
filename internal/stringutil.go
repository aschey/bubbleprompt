package internal

import (
	"strings"
)

func AddNewlineIfMissing(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s
	} else {
		return s + "\n"
	}
}

func TrimNewline(s string) string {
	return strings.TrimSuffix(s, "\n")
}

func CountNewlines(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	return count
}

func Unescape(s string, wrapper string) string {
	val := s
	if strings.HasPrefix(val, wrapper) {
		val = strings.TrimPrefix(val, wrapper)
		val = strings.TrimSuffix(val, wrapper)
	}
	val = strings.ReplaceAll(val, `\`+wrapper, wrapper)

	return val
}
