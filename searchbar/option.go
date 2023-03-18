package searchbar

import (
	prompt "github.com/aschey/bubbleprompt"
	"github.com/charmbracelet/lipgloss"
)

type searchbarSettings[T any] struct {
	searchbarStyle lipgloss.Style
	label          string
	maxWidth       int
	promptOptions  []prompt.Option[T]
}

type Option[T any] func(settings *searchbarSettings[T])

func WithSearchbarStyle[T any](searchbarStyle lipgloss.Style) Option[T] {
	return func(settings *searchbarSettings[T]) {
		settings.searchbarStyle = searchbarStyle
	}
}

func WithLabel[T any](label string) Option[T] {
	return func(settings *searchbarSettings[T]) {
		settings.label = label
	}
}

func WithMaxWidth[T any](maxWidth int) Option[T] {
	return func(settings *searchbarSettings[T]) {
		settings.maxWidth = maxWidth
	}
}

func WithPromptOptions[T any](options ...prompt.Option[T]) Option[T] {
	return func(settings *searchbarSettings[T]) {
		settings.promptOptions = options
	}
}
