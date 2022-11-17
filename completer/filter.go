package completer

import (
	"strings"

	"github.com/aschey/bubbleprompt/editor"
)

func FilterHasPrefix[T any](search string, suggestions []editor.Suggestion[T]) []editor.Suggestion[T] {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))

	filtered := []editor.Suggestion[T]{}
	for _, s := range suggestions {
		completionText := strings.ToLower(s.GetCompletionText())
		if strings.HasPrefix(completionText, cleanedSearch) || strings.HasPrefix(s.Text, cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
