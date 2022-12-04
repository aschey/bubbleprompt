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

type CompleteMsg struct{}

type SuggestionMsg[T any] struct {
	Suggestions []Suggestion[T]
	Err         error
}

func Complete() tea.Msg {
	return CompleteMsg{}
}
