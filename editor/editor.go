package editor

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewMode int

const (
	Interactive ViewMode = iota
	Static
)

type Editor[T any] interface {
	Init() tea.Cmd
	OnUpdateStart(msg tea.Msg) tea.Cmd
	View(viewMode ViewMode) string
	Focus() tea.Cmd
	Focused() bool
	Value() string
	Runes() []rune
	ResetValue()
	SetValue(value string)
	Blur()
	CursorIndex() int
	CursorOffset() int
	SetCursor(cursor int)
	SetCursorMode(cursorMode textinput.CursorMode) tea.Cmd
	Prompt() string
	SetPrompt(prompt string)
	Tokens() []Token
	ShouldSelectSuggestion(suggestion Suggestion[T]) bool
	SuggestionRunes(runes []rune) []rune
	OnUpdateFinish(msg tea.Msg, suggestion *Suggestion[T], isSelected bool) tea.Cmd
	OnSuggestionChanged(suggestion Suggestion[T])
	OnExecutorFinished()
	OnSuggestionUnselected()
	ShouldClearSuggestions(prevRunes []rune, msg tea.KeyMsg) bool
	ShouldUnselectSuggestion(prevRunes []rune, msg tea.KeyMsg) bool
}
