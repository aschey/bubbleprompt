package completer

import (
	"github.com/aschey/bubbleprompt/editor"
)

type Metadata[T any] interface {
	Children() []editor.Suggestion[T]
}

func GetRecursiveSuggestions[T Metadata[T]](tokens []editor.Token, cursor int, suggestions []editor.Suggestion[T]) []editor.Suggestion[T] {
	if len(tokens) == 0 {
		return suggestions
	}
	token := tokens[0]

	if cursor <= token.End() {
		prefixEnd := cursor - token.Start
		if prefixEnd < 0 {
			return []editor.Suggestion[T]{}
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
			return []editor.Suggestion[T]{}
		}
	}
	return []editor.Suggestion[T]{}
}
