package main

import (
	"fmt"
	"os"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	tea "github.com/charmbracelet/bubbletea"
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

func (m completerModel) completer(input string) prompt.Suggestions {
	time.Sleep(100 * time.Millisecond)
	return prompt.FilterHasPrefix(input, m.suggestions)
}

func executor(input string, selected *prompt.Suggestion, suggestions prompt.Suggestions) tea.Model {
	return prompt.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "test"
	})
}

func main() {
	suggestions := []prompt.Suggestion{
		{Name: "first option", Description: "test desc", Placeholder: "[hh]"},
		{Name: "second option", Description: "test desc2"},
		{Name: "third option", Description: "test desc2"},
		{Name: "fourth option", Description: "test desc2"},
		{Name: "fifth option", Description: "test desc2"},
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
