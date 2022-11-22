package completer

import "github.com/aschey/bubbleprompt/input"

type Metadata[T any] interface {
	Children() []input.Suggestion[T]
}

func GetRecursiveSuggestions[T Metadata[T]](tokens []input.Token, cursor int, suggestions []input.Suggestion[T]) []input.Suggestion[T] {
	if len(tokens) == 0 {
		return suggestions
	}
	token := tokens[0]

	if cursor <= token.End() {
		prefixEnd := cursor - token.Start
		if prefixEnd < 0 {
			return []input.Suggestion[T]{}
		}
		return FilterHasPrefix(string([]rune(token.Value)[:cursor-token.Start]), suggestions)
	}
	for _, suggestion := range suggestions {
		if suggestion.GetSuggestionText() == token.Value {
			if suggestion.Metadata.Children() != nil {
				children := suggestion.Metadata.Children()
				if children != nil {
					return GetRecursiveSuggestions(tokens[1:], cursor, children)
				}
			}
			return []input.Suggestion[T]{}
		}
	}
	return []input.Suggestion[T]{}
}
