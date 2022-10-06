package completer

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
)

func FilterHasPrefix[T any](search string, suggestions []input.Suggestion[T]) []input.Suggestion[T] {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))
	filtered := []input.Suggestion[T]{}
	for _, s := range suggestions {
		if strings.HasPrefix(s.GetCompletionText(), cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
