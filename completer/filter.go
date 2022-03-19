package completer

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
)

func FilterHasPrefix(search string, suggestions []input.Suggestion) []input.Suggestion {
	return filterHasPrefix(search, suggestions,
		func(s input.Suggestion) string { return s.Text })
}

func FilterCompletionTextHasPrefix(search string, suggestions []input.Suggestion) []input.Suggestion {
	return filterHasPrefix(search, suggestions,
		func(s input.Suggestion) string { return s.CompletionText })
}

func filterHasPrefix(search string, suggestions []input.Suggestion,
	textFunc func(s input.Suggestion) string) []input.Suggestion {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))
	filtered := []input.Suggestion{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(textFunc(s)), cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
