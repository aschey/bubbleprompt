package tutorial

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/charmbracelet/lipgloss"
)

func main3() {
	// Initialize the input
	textInput := simpleinput.New[any]()

	// Define our suggestions
	suggestions := []input.Suggestion[any]{
		{Text: "banana", Description: "good with peanut butter"},
		{Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberry", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}

	model := model{
		suggestions: suggestions,
		textInput:   textInput,
		// Add some coloring to the foreground of our output to make it look pretty
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}
	_ = model
}
