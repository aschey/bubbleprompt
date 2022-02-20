package prompt

import tea "github.com/charmbracelet/bubbletea"

type StringModel string

func NewStringModel(output string) StringModel {
	return StringModel(output)
}

func (s StringModel) Init() tea.Cmd {
	return nil
}

func (s StringModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return s, tea.Quit
}

func (s StringModel) View() string {
	return string(s)
}
