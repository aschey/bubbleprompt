package tutorial

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main4() {
	// Initialize the input
	textInput := simpleinput.New[any]()

	// Define our suggestions
	suggestions := []suggestion.Suggestion[any]{
		{Text: "banana", Description: "good with peanut butter"},
		{Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberry", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}

	// Combine everything into our model
	model := model{
		suggestions: suggestions,
		textInput:   textInput,
		// Add some coloring to the foreground of our output to make it look pretty
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}

	// Create the Bubbleprompt model
	// This struct fulfills the tea.Model interface so it can be passed directly to tea.NewProgram
	promptModel := prompt.New[any](model, textInput)

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
	fmt.Println()

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
