package input

import (
	tea "github.com/charmbracelet/bubbletea"
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
	Prompt() string
	SetPrompt(string)
	CompletionText(string) string
	OnUpdateFinish(msg tea.Msg, suggestion *Suggestion) tea.Cmd
	OnSuggestionChanged(suggestion Suggestion)
	IsDelimiter(string) bool
	OnSuggestionUnselected()
}
