package completer

import (
	"github.com/aschey/bubbleprompt/input"
)

type Metadata[T any] interface {
	Children() []input.Suggestion[T]
}

func GetRecursiveCompletions[T Metadata[T]](tokens []input.Token, cursor int, completions []input.Suggestion[T]) []input.Suggestion[T] {
	if len(tokens) == 0 {
		return completions
	}
	token := tokens[0]

	if cursor <= token.End() {
		prefixEnd := cursor - token.Start
		if prefixEnd < 0 {
			return []input.Suggestion[T]{}
		}
		return FilterHasPrefix(string([]rune(token.Value)[:cursor-token.Start]), completions)
	}
	for _, completion := range completions {
		if completion.GetCompletionText() == token.Value {
			if completion.Metadata.Children() != nil {
				children := completion.Metadata.Children()
				if children != nil {
					return GetRecursiveCompletions(tokens[1:], cursor, children)
				}
			}
			return []input.Suggestion[T]{}
		}
	}
	return []input.Suggestion[T]{}
}
