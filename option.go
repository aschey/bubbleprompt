package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/charmbracelet/lipgloss"
)

type Option[T any] func(model *Model[T]) error

func WithFormatters[T any](formatters input.Formatters) Option[T] {
	return func(model *Model[T]) error {
		model.Formatters = formatters
		return nil
	}
}

func WithScrollbar[T any](color lipgloss.TerminalColor) Option[T] {
	return func(model *Model[T]) error {
		model.SetScrollbarColor(color)
		return nil
	}
}

func WithScrollbarThumb[T any](color lipgloss.TerminalColor) Option[T] {
	return func(model *Model[T]) error {
		model.SetScrollbarThumbColor(color)
		return nil
	}
}

func WithMaxSuggestions[T any](maxSuggestions int) Option[T] {
	return func(model *Model[T]) error {
		model.SetMaxSuggestions(maxSuggestions)
		return nil
	}
}

func WithUnmanagedRenderer[T any]() Option[T] {
	return func(model *Model[T]) error {
		model.SetRenderer(renderer.NewUnmanagedRenderer())
		return nil
	}
}

func WithViewportRenderer[T any]() Option[T] {
	return func(model *Model[T]) error {
		model.SetRenderer(renderer.NewViewportRenderer())
		return nil
	}
}
