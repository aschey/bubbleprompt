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
	SetContent(content string)
	History() string
	AddOutput(output string)
	GotoBottom(msg tea.Msg)
}
