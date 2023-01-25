package completer

import (
	"strings"

	"github.com/aschey/bubbleprompt/suggestion"
)

func FilterHasPrefix[T any](
	search string,
	suggestions []suggestion.Suggestion[T],
) []suggestion.Suggestion[T] {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))

	filtered := []suggestion.Suggestion[T]{}
	for _, s := range suggestions {
		suggestionText := strings.ToLower(s.GetSuggestionText())
		if strings.HasPrefix(suggestionText, cleanedSearch) ||
			strings.HasPrefix(s.Text, cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
