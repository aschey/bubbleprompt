package completer

import (
	"github.com/aschey/bubbleprompt/editor"
)

type Metadata[T any] interface {
	Children() []editor.Suggestion[T]
}

func GetRecursiveCompletions[T Metadata[T]](tokens []editor.Token, cursor int, completions []editor.Suggestion[T]) []editor.Suggestion[T] {
	if len(tokens) == 0 {
		return completions
	}
	token := tokens[0]

	if cursor <= token.End() {
		prefixEnd := cursor - token.Start
		if prefixEnd < 0 {
			return []editor.Suggestion[T]{}
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
			return []editor.Suggestion[T]{}
		}
	}
	return []editor.Suggestion[T]{}
}
