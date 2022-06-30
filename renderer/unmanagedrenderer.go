package renderer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UnmanagedRenderer struct {
	content string
	history string
}

func NewUnmanagedRenderer() *UnmanagedRenderer {
	return &UnmanagedRenderer{}
}

func (u *UnmanagedRenderer) View() string {
	return u.content
}

func (u *UnmanagedRenderer) Initialize(msg tea.WindowSizeMsg) {}

func (u *UnmanagedRenderer) SetSize(msg tea.WindowSizeMsg) {}

func (u *UnmanagedRenderer) Update(msg tea.Msg) (Renderer, tea.Cmd) {
	return u, nil
}

func (u *UnmanagedRenderer) SetContent(content string) {
	u.content = content
}

func (u *UnmanagedRenderer) History() string {
	return ""
}

func (u *UnmanagedRenderer) AddOutput(output string) {
	if len(u.history) > 0 {
		u.history += "\n"
	}
	u.history += output
}

func (u *UnmanagedRenderer) GotoBottom(msg tea.Msg) {}

func (u *UnmanagedRenderer) FinishUpdate() tea.Cmd {
	if len(u.history) == 0 {
		return nil
	}
	history := u.history
	u.history = ""

	return tea.Println(history)
}
