package completer

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/suggestion"
)

type Metadata[T any] interface {
	GetChildren() []suggestion.Suggestion[T]
}

type RecursiveFilterer[T Metadata[T]] struct {
	Filterer Filterer[T]
}

func NewRecursiveFilterer[T Metadata[T]]() RecursiveFilterer[T] {
	return RecursiveFilterer[T]{}
}

func (f RecursiveFilterer[T]) GetRecursiveSuggestions(
	tokens []input.Token,
	cursor int,
	suggestions []suggestion.Suggestion[T],
) []suggestion.Suggestion[T] {
	if len(tokens) == 0 {
		return suggestions
	}
	token := tokens[0]

	if cursor <= token.End() {
		prefixEnd := cursor - token.Start
		if prefixEnd < 0 {
			return []suggestion.Suggestion[T]{}
		}

		return f.getFilterer().Filter(string([]rune(token.Value)[:cursor-token.Start]), suggestions)
	}
	for _, sug := range suggestions {
		if sug.GetSuggestionText() == token.Value {
			if sug.Metadata.GetChildren() != nil {
				children := sug.Metadata.GetChildren()
				if children != nil {
					return f.GetRecursiveSuggestions(tokens[1:], cursor, children)
				}
			}
			return []suggestion.Suggestion[T]{}
		}
	}
	return []suggestion.Suggestion[T]{}
}

func (s *RecursiveFilterer[T]) getFilterer() Filterer[T] {
	if s.Filterer == nil {
		s.Filterer = NewPrefixFilter[T]()
	}
	return s.Filterer
}
