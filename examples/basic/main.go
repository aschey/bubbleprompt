package main

import (
	"fmt"
	"os"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	prompt prompt.Model
}

type completerModel struct {
	suggestions []prompt.Suggestion
}

func (m model) Init() tea.Cmd {
	return m.prompt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.prompt.Update(msg)
	m.prompt = p
	return m, cmd
}

func (m model) View() string {
	return m.prompt.View()
}

func (m completerModel) completer(document prompt.Document) prompt.Suggestions {
	time.Sleep(100 * time.Millisecond)
	return prompt.FilterHasPrefix(document.InputBeforeCursor, m.suggestions)
}

func executor(input string, selected *prompt.Suggestion, suggestions prompt.Suggestions) tea.Model {
	return prompt.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "test"
	})
}

func main() {
	placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	argStyle1 := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	argStyle2 := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	suggestions := []prompt.Suggestion{
		{Name: "first-option", Description: "test description",
			PositionalArgs: []prompt.PositionalArg{
				{Placeholder: "test1", PlaceholderStyle: prompt.Text{Style: placeholderStyle}, ArgStyle: prompt.Text{Style: argStyle1}},
				{Placeholder: "test2", PlaceholderStyle: prompt.Text{Style: placeholderStyle}, ArgStyle: prompt.Text{Style: argStyle2}},
			}},
		{Name: "second-option", Description: "test description2"},
		{Name: "third-option", Description: "test description3"},
		{Name: "fourth-option", Description: "test description4"},
		{Name: "fifth-option", Description: "test description5",
			PositionalArgs: []prompt.PositionalArg{
				{Placeholder: "abc", PlaceholderStyle: prompt.Text{Style: placeholderStyle}},
			}},
	}

	completerModel := completerModel{suggestions: suggestions}

	m := model{prompt: prompt.New(
		completerModel.completer,
		executor,
		prompt.WithPrompt(">>> "),
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
