package main

import (
	"fmt"
	"os"
	"strconv"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	suggestions []suggestion.Suggestion[any]
	textInput   *simpleinput.Model[any]
	outputStyle lipgloss.Style
	numChoices  int64
	filterer    completer.Filterer[any]
}

func (m model) Complete(promptModel prompt.Model[any]) ([]suggestion.Suggestion[any], error) {
	if len(m.textInput.Tokens()) > 1 {
		return nil, nil
	}

	return m.filterer.Filter(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	tokens := m.textInput.WordTokenValues()
	if len(tokens) == 0 {
		return nil, fmt.Errorf("No selection")
	}
	return executor.NewStringModel(m.formatOutput(tokens[0])), nil
}

func (m model) formatOutput(choice string) string {
	return fmt.Sprintf("You picked: %s\nYou've entered %s submissions(s)\n\n",
		m.outputStyle.Render(choice),
		m.outputStyle.Render(strconv.FormatInt(m.numChoices, 10)))
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		m.numChoices++
	}
	return m, nil
}

func main() {
	textInput := simpleinput.New[any]()
	suggestions := []suggestion.Suggestion[any]{
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
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
		filterer:    completer.NewPrefixFilter[any](),
	}

	promptModel := prompt.New[any](model, textInput)

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
	fmt.Println()

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
