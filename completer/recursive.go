package completer

import (
	"github.com/aschey/bubbleprompt/input"
)

type Node[T any] interface {
	pushNext(new Node[T])
	valid() bool
	next() Node[T]
}

type Metadata interface {
	Children() []input.Suggestion[Metadata]
}

func GetRecursiveCompletions(tokens []input.Token, cursor int, completions []input.Suggestion[Metadata]) []input.Suggestion[Metadata] {
	if len(tokens) == 0 {
		return completions
	}
	token := tokens[0]

	if cursor <= token.End() {
		prefixEnd := cursor - token.Start
		if prefixEnd < 0 {
			return []input.Suggestion[Metadata]{}
		}
		return FilterHasPrefix(string([]rune(token.Value)[:cursor-token.Start]), completions)
	}
	for _, completion := range completions {
		if completion.GetCompletionText() == token.Value {
			if completion.Metadata != nil {
				children := completion.Metadata.Children()
				if children != nil {
					return GetRecursiveCompletions(tokens[1:], cursor, children)
				}
			}
			return []input.Suggestion[Metadata]{}
		}
	}
	return []input.Suggestion[Metadata]{}
}
