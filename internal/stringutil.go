package internal

import "strings"

func AddNewlineIfMissing(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s
	} else {
		return s + "\n"
	}
}

func TrimNewline(s string) string {
	return strings.TrimRight(s, "\n")
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
