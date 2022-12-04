package tutorial

import (
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	// list of suggestions that we'll display using the completer function
	suggestions []suggestion.Suggestion[any]
	// Reference to our input component. We'll use this to read user input
	textInput *simpleinput.Model[any]
	// Style struct for formatting the output
	outputStyle lipgloss.Style
	// Number of times the user enters some input
	numChoices int64
}
