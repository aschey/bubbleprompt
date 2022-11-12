package completer

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
)

func FilterHasPrefix[T any](search string, suggestions []input.Suggestion[T]) []input.Suggestion[T] {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))

	filtered := []input.Suggestion[T]{}
	for _, s := range suggestions {
		completionText := strings.ToLower(s.GetCompletionText())
		if strings.HasPrefix(completionText, cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
