package suggestion

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type PeriodicCompleterMsg struct {
	NextTrigger time.Duration
}

type OneShotCompleterMsg struct{}

func PeriodicCompleter(nextTrigger time.Duration) tea.Cmd {
	return tea.Tick(nextTrigger, func(_ time.Time) tea.Msg {
		return PeriodicCompleterMsg{NextTrigger: nextTrigger}
	})
}

func OneShotCompleter(nextTrigger time.Duration) tea.Cmd {
	return tea.Tick(nextTrigger, func(_ time.Time) tea.Msg {
		return OneShotCompleterMsg{}
	})
}

type RefreshSuggestionsMessage[T any] []Suggestion[T]

func RefreshSuggestions[T any](init func() []Suggestion[T]) tea.Cmd {
	return func() tea.Msg {
		return RefreshSuggestionsMessage[T](init())
	}
}

type CompleteMsg struct{}

type SuggestionMsg[T any] struct {
	Suggestions    []Suggestion[T]
	SequenceNumber int
	Err            error
}

func Complete() tea.Msg {
	return CompleteMsg{}
}
