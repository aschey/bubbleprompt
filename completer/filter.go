package completer

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
)

func FilterHasPrefix[T any](search string, suggestions []input.Suggestion[T]) []input.Suggestion[T] {
	return filterHasPrefix(search, suggestions,
		func(s input.Suggestion[T]) string { return s.Text })
}

func FilterCompletionTextHasPrefix[T any](search string, suggestions []input.Suggestion[T]) []input.Suggestion[T] {
	return filterHasPrefix(search, suggestions,
		func(s input.Suggestion[T]) string { return s.CompletionText })
}

func filterHasPrefix[T any](search string, suggestions []input.Suggestion[T],
	textFunc func(s input.Suggestion[T]) string) []input.Suggestion[T] {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))
	filtered := []input.Suggestion[T]{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(textFunc(s)), cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
