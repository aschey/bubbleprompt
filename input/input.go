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
	SetValue(value string)
	Blur()
	Cursor() int
	SetCursor(cursor int)
	Prompt() string
	SetPrompt(prompt string)
	CompletionText(text string) string
	OnUpdateFinish(msg tea.Msg, suggestion *Suggestion) tea.Cmd
	OnSuggestionChanged(suggestion Suggestion)
	IsDelimiter(text string) bool
	OnSuggestionUnselected()
	ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool
}
