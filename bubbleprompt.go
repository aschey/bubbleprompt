package bubbleprompt

import (
	tea "github.com/charmbracelet/bubbletea"
)

type suggest struct {
	name        string
	description string
}

type Model struct {
	suggestions []suggest
}

func New() Model {
	return Model{suggestions: []suggest{
		{name: "test name", description: "test desc"},
		{name: "test name2", description: "test desc2"},
	}}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ">>>"
}
