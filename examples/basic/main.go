package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	promptModel prompt.Model[any]
}

type completerModel struct {
	suggestions []input.Suggestion[any]
	textInput   *simpleinput.Model
	outputStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return m.promptModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.promptModel.Update(msg)
	m.promptModel = p
	return m, cmd
}

func (m model) View() string {
	return m.promptModel.View()
}

func (m completerModel) completer(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
	if len(m.textInput.AllTokens()) > 1 {
		return nil, nil
	}

	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[any]) (tea.Model, error) {
	return executors.NewStringModel(m.outputStyle.Render("You picked: " + input)), nil
}

func main() {
	textInput := simpleinput.New()
	suggestions := []input.Suggestion[any]{
		{Text: "apple", Description: "spherical...ish"},
		{Text: "banana", Description: "good with peanut butter"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberry", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}

	completerModel := completerModel{
		suggestions: suggestions,
		textInput:   textInput,
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
	}

	promptModel, err := prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	)
	if err != nil {
		panic(err)
	}

	m := model{promptModel}
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
	fmt.Println()
	if _, err := tea.NewProgram(m, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
