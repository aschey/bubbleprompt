package renderer

import (
	"github.com/aschey/bubbleprompt/internal"
	tea "github.com/charmbracelet/bubbletea"
)

type UnmanagedRenderer struct {
	content        string
	currentHistory string
	totalHistory   string
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

func (u *UnmanagedRenderer) SetHistory(history string) tea.Cmd {
	u.totalHistory = history
	if len(history) > 0 {
		return tea.Println(internal.TrimNewline(history))
	}
	return nil
}

func (u *UnmanagedRenderer) GetHistory() string {
	return u.totalHistory
}

func (u *UnmanagedRenderer) AddOutput(output string) {
	if len(u.currentHistory) > 0 {
		u.currentHistory = internal.AddNewlineIfMissing(u.currentHistory)
	}
	if len(u.totalHistory) > 0 {
		u.totalHistory = internal.AddNewlineIfMissing(u.totalHistory)
	}

	u.currentHistory += internal.TrimNewline(output)
	u.totalHistory += internal.TrimNewline(output)
}

func (u *UnmanagedRenderer) GotoBottom(msg tea.Msg) {}

func (u *UnmanagedRenderer) FinishUpdate() tea.Cmd {
	if len(u.currentHistory) == 0 {
		return nil
	}
	currentHistory := u.currentHistory
	u.currentHistory = ""

	return tea.Println(currentHistory)
}
