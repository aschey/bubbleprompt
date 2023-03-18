package suggestion

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Manager[T any] interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	SetMaxSuggestions(maxSuggestions int)
	MaxSuggestions() int
	SetSelectionIndicator(selectionIndicator string)
	SelectionIndicator() string
	EnableScrollbar()
	DisableScrollbar()
	SelectedSuggestion() *Suggestion[T]
	SelectedIndex() int
	Suggestions() []Suggestion[T]
	VisibleSuggestions() []Suggestion[T]
	MaxSuggestionWidth() (int, int)
	SelectSuggestion(suggestion Suggestion[T])
	UnselectSuggestion()
	IsSuggestionSelected() bool
	PreviousSuggestion()
	NextSuggestion()
	UpdateSuggestions() tea.Cmd
	ResetSuggestions() tea.Cmd
	ClearSuggestions()
	Error() error
	ScrollbarBounds() (int, int)
	ScrollPosition() int
	Scrollbar() string
	ScrollbarThumb() string
	Render(paddingSize int) string
	ShouldChangeListPosition(msg tea.Msg) bool
	Formatters() Formatters
	SetFormatters(formatters Formatters)
}
