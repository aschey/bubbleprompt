package renderer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Renderer interface {
	View() string
	Initialize(msg tea.WindowSizeMsg)
	SetSize(msg tea.WindowSizeMsg)
	Update(msg tea.Msg) (Renderer, tea.Cmd)
	FinishUpdate() tea.Cmd
	SetInput(input string)
	SetBody(suggestions string)
	AddHistory(output string)
	GotoBottom(msg tea.Msg)
	GetHistory() string
	SetHistory(history string) tea.Cmd
}
