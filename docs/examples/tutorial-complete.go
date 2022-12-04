package tutorial

import (
	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/suggestion"
)

func (m model) Complete(promptModel prompt.Model[any]) ([]suggestion.Suggestion[any], error) {
	// Our program only takes one token as input,
	// so don't return any suggestions if the user types more than one word
	if len(m.textInput.Tokens()) > 1 {
		return nil, nil
	}

	// Filter suggestions based on the text before the cursor
	return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}
