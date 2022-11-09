package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cmdMetadata = commandinput.CmdMetadata

type model struct {
	promptModel prompt.Model[cmdMetadata]
}

type completerModel struct {
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
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

func (m completerModel) completer(promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	if m.textInput.CommandCompleted() {
		return nil, nil
	}

	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
	return executors.NewStringModel(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("You picked: " + input)), nil
}

func main() {
	textInput := commandinput.New[cmdMetadata]()
	suggestions := []input.Suggestion[cmdMetadata]{
		{Text: "apples", Description: "spherical...ish"},
		{Text: "bananas", Description: "good with peanut butter"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberries", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}

	completerModel := completerModel{suggestions: suggestions, textInput: textInput}

	promptModel, err := prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	)
	if err != nil {
		panic(err)
	}
	m := model{promptModel}
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("pick a fruit!"))
	fmt.Println()
	if _, err := tea.NewProgram(m, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
