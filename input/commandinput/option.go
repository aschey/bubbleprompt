package commandinput

import "github.com/charmbracelet/bubbles/cursor"

type Option[T any] func(model *Model[T])

func WithPrompt[T any](prompt string) Option[T] {
	return func(model *Model[T]) {
		model.SetPrompt(prompt)
	}
}

func WithFormatters[T any](formatters Formatters) Option[T] {
	return func(model *Model[T]) {
		model.SetFormatters(formatters)
	}
}

func WithDefaultDelimiter[T any](defaultDelimiter string) Option[T] {
	return func(model *Model[T]) {
		model.defaultDelimiter = defaultDelimiter
	}
}

func WithCursorMode[T any](cursorMode cursor.Mode) Option[T] {
	return func(model *Model[T]) {
		model.SetCursorMode(cursorMode)
	}
}
