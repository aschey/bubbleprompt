package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input interface {
	Init() tea.Cmd
	OnUpdateStart(msg tea.Msg) tea.Cmd
	View() string
	Focus() tea.Cmd
	Focused() bool
	Value() string
	SetValue(string)
	Blur()
	Cursor() int
	SetCursor(int)
	SetTextStyle(lipgloss.Style)
	Prompt() string
	SetPrompt(string)
	CompletionText(string) string
	OnUpdateFinish(msg tea.Msg, suggestion *Suggestion) tea.Cmd
	OnSuggestionChanged(suggestion Suggestion)
	IsDelimiter(string) bool
}
