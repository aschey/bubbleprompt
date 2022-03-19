package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/lipgloss"
)

type Option func(model *Model) error

func WithFormatters(formatters input.Formatters) Option {
	return func(model *Model) error {
		model.Formatters = formatters
		return nil
	}
}

func WithScrollbar(color lipgloss.TerminalColor) Option {
	return func(model *Model) error {
		model.SetScrollbarColor(color)
		return nil
	}
}

func WithScrollbarThumb(color lipgloss.TerminalColor) Option {
	return func(model *Model) error {
		model.SetScrollbarThumbColor(color)
		return nil
	}
}

func WithMaxSuggestions(maxSuggestions int) Option {
	return func(model *Model) error {
		model.SetMaxSuggestions(maxSuggestions)
		return nil
	}
}
