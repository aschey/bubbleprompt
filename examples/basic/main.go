package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt prompt.Model
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

func main() {
	suggestions := []prompt.Suggest{
		{Name: "first option", Description: "test desc"},
		{Name: "second option", Description: "test desc2"},
		{Name: "third option", Description: "test desc2"},
		{Name: "fourth option", Description: "test desc2"},
		{Name: "fifth option", Description: "test desc2"},
	}
	m := model{prompt: prompt.New(prompt.OptionInitialSuggestions(suggestions), prompt.OptionPrompt(">>> "))}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
