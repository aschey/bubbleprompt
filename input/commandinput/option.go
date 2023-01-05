package commandinput

import "github.com/charmbracelet/bubbles/cursor"

type Option[T CommandMetadataAccessor] func(model *Model[T])

func WithPrompt[T CommandMetadataAccessor](prompt string) Option[T] {
	return func(model *Model[T]) {
		model.SetPrompt(prompt)
	}
}

func WithFormatters[T CommandMetadataAccessor](formatters Formatters) Option[T] {
	return func(model *Model[T]) {
		model.SetFormatters(formatters)
	}
}

func WithDefaultDelimiter[T CommandMetadataAccessor](defaultDelimiter string) Option[T] {
	return func(model *Model[T]) {
		model.defaultDelimiter = defaultDelimiter
	}
}

func WithCursorMode[T CommandMetadataAccessor](cursorMode cursor.Mode) Option[T] {
	return func(model *Model[T]) {
		model.SetCursorMode(cursorMode)
	}
}
