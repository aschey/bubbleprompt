package completer

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/suggestion"
)

type Metadata[T any] interface {
	Children() []suggestion.Suggestion[T]
}

func GetRecursiveSuggestions[T Metadata[T]](tokens []input.Token, cursor int, suggestions []suggestion.Suggestion[T]) []suggestion.Suggestion[T] {
	if len(tokens) == 0 {
		return suggestions
	}
	token := tokens[0]

	if cursor <= token.End() {
		prefixEnd := cursor - token.Start
		if prefixEnd < 0 {
			return []suggestion.Suggestion[T]{}
		}
		return FilterHasPrefix(string([]rune(token.Value)[:cursor-token.Start]), suggestions)
	}
	for _, s := range suggestions {
		if s.GetSuggestionText() == token.Value {
			if s.Metadata.Children() != nil {
				children := s.Metadata.Children()
				if children != nil {
					return GetRecursiveSuggestions(tokens[1:], cursor, children)
				}
			}
			return []suggestion.Suggestion[T]{}
		}
	}
	return []suggestion.Suggestion[T]{}
}
