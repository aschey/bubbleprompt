package main

import (
	"fmt"
	"os"

	bubbleprompt "github.com/aschey/bubbleprompt"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt bubbleprompt.Model
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
	m := model{prompt: bubbleprompt.New()}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
