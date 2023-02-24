package completer

import (
	"strings"

	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/sahilm/fuzzy"
)

type Filterer[T any] interface {
	Filter(
		search string,
		suggestions []suggestion.Suggestion[T],
	) []suggestion.Suggestion[T]
}

type PrefixFilter[T any] struct{}

func NewPrefixFilter[T any]() PrefixFilter[T] {
	return PrefixFilter[T]{}
}

func (f PrefixFilter[T]) Filter(
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

type FuzzyFilter[T any] struct{}

func NewFuzzyFilter[T any]() FuzzyFilter[T] {
	return FuzzyFilter[T]{}
}

type suggestionSource[T any] []suggestion.Suggestion[T]

func (s suggestionSource[T]) String(i int) string {
	return s[i].GetSuggestionText()
}

func (s suggestionSource[T]) Len() int {
	return len(s)
}

func (f FuzzyFilter[T]) Filter(
	search string,
	suggestions []suggestion.Suggestion[T],
) []suggestion.Suggestion[T] {
	if search == "" {
		return suggestions
	}

	matches := fuzzy.FindFrom(search, suggestionSource[T](suggestions))
	filtered := []suggestion.Suggestion[T]{}
	for _, match := range matches {
		filtered = append(filtered, suggestions[match.Index])
	}

	return filtered
}
