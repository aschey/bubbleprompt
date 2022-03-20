package executor

import tea "github.com/charmbracelet/bubbletea"

type StringModel string

func NewStringModel(output string) StringModel {
	return StringModel(output)
}

func (m StringModel) Init() tea.Cmd {
	return nil
}

func (m StringModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m StringModel) View() string {
	return string(m + "\n")
}
