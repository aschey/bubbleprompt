package input

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Input[T any] interface {
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
	ShouldSelectSuggestion(suggestion Suggestion[T]) bool
	CompletionText(text string) string
	OnUpdateFinish(msg tea.Msg, suggestion *Suggestion[T]) tea.Cmd
	OnSuggestionChanged(suggestion Suggestion[T])
	OnExecutorFinished()
	OnSuggestionUnselected()
	ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool
	ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool
}
