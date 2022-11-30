package commandinput

import "github.com/charmbracelet/bubbles/cursor"

type Option[T CmdMetadataAccessor] func(model *Model[T])

func WithPrompt[T CmdMetadataAccessor](prompt string) Option[T] {
	return func(model *Model[T]) {
		model.SetPrompt(prompt)
	}
}

func WithFormatters[T CmdMetadataAccessor](formatters Formatters) Option[T] {
	return func(model *Model[T]) {
		model.SetFormatters(formatters)
	}
}

func WithDefaultDelimiter[T CmdMetadataAccessor](defaultDelimiter string) Option[T] {
	return func(model *Model[T]) {
		model.defaultDelimiter = defaultDelimiter
	}
}

func WithCursorMode[T CmdMetadataAccessor](cursorMode cursor.Mode) Option[T] {
	return func(model *Model[T]) {
		model.SetCursorMode(cursorMode)
	}
}
