package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type completerModel struct {
	suggestions []input.Suggestion[any]
	textInput   *simpleinput.Model[any]
	outputStyle lipgloss.Style
}

func (m completerModel) completer(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
	if len(m.textInput.AllTokens()) > 1 {
		return nil, nil
	}

	return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[any]) (tea.Model, error) {
	tokens := m.textInput.TokenValues()
	if len(tokens) == 0 {
		return executor.NewStringModel("No selection"), nil
	}
	return executor.NewStringModel("You picked: " + m.outputStyle.Render(tokens[0])), nil
}

func main() {
	textInput := simpleinput.New[any]()
	suggestions := []input.Suggestion[any]{
		{Text: "banana", Description: "good with peanut butter"},
		{Text: "\"sugar apple\"", CompletionText: "sugar apple", Description: "spherical...ish"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberry", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}

	completerModel := completerModel{
		suggestions: suggestions,
		textInput:   textInput,
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}

	promptModel, err := prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
	fmt.Println()
	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
