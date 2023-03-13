package renderer

import (
	"github.com/aschey/bubbleprompt/internal"
	tea "github.com/charmbracelet/bubbletea"
)

type UnmanagedRenderer struct {
	content        string
	currentHistory string
	totalHistory   string
	settings       rendererSettings
}

func NewUnmanagedRenderer(options ...Option) *UnmanagedRenderer {
	settings := rendererSettings{
		widthOffset:  0,
		heightOffset: 0,
		useHistory:   true,
	}
	for _, option := range options {
		option(&settings)
	}
	return &UnmanagedRenderer{settings: settings}
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
	if u.settings.useHistory {
		u.totalHistory = history
		if len(history) > 0 {
			return tea.Println(internal.TrimNewline(history))
		}
	}
	return nil
}

func (u *UnmanagedRenderer) GetHistory() string {
	return u.totalHistory
}

func (u *UnmanagedRenderer) AddOutput(output string) {
	if u.settings.useHistory {
		if len(u.currentHistory) > 0 {
			u.currentHistory = internal.AddNewlineIfMissing(u.currentHistory)
		}
		if len(u.totalHistory) > 0 {
			u.totalHistory = internal.AddNewlineIfMissing(u.totalHistory)
		}

		u.currentHistory += internal.TrimNewline(output)
		u.totalHistory += internal.TrimNewline(output)
	}
}

func (u *UnmanagedRenderer) GotoBottom(msg tea.Msg) {}

func (u *UnmanagedRenderer) FinishUpdate() tea.Cmd {
	if u.settings.useHistory {
		if len(u.currentHistory) == 0 {
			return nil
		}
		currentHistory := u.currentHistory
		u.currentHistory = ""

		return tea.Println(currentHistory)
	}
	return nil
}
