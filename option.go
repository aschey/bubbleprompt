package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/lipgloss"
)

type Option[I any] func(model *Model[I]) error

func WithFormatters[I any](formatters input.Formatters) Option[I] {
	return func(model *Model[I]) error {
		model.Formatters = formatters
		return nil
	}
}

func WithScrollbar[I any](color lipgloss.TerminalColor) Option[I] {
	return func(model *Model[I]) error {
		model.SetScrollbarColor(color)
		return nil
	}
}

func WithScrollbarThumb[I any](color lipgloss.TerminalColor) Option[I] {
	return func(model *Model[I]) error {
		model.SetScrollbarThumbColor(color)
		return nil
	}
}

func WithMaxSuggestions[I any](maxSuggestions int) Option[I] {
	return func(model *Model[I]) error {
		model.SetMaxSuggestions(maxSuggestions)
		return nil
	}
}
