package commandinput

import "github.com/charmbracelet/bubbles/cursor"

type Option[T CmdMetadataAccessor] func(model *Model[T]) error

func WithPrompt[T CmdMetadataAccessor](prompt string) Option[T] {
	return func(model *Model[T]) error {
		model.SetPrompt(prompt)
		return nil
	}
}

func WithFormatters[T CmdMetadataAccessor](formatters Formatters) Option[T] {
	return func(model *Model[T]) error {
		model.SetFormatters(formatters)
		return nil
	}
}

func WithDefaultDelimiter[T CmdMetadataAccessor](defaultDelimiter string) Option[T] {
	return func(model *Model[T]) error {
		model.defaultDelimiter = defaultDelimiter
		return nil
	}
}

func WithCursorMode[T CmdMetadataAccessor](cursorMode cursor.Mode) Option[T] {
	return func(model *Model[T]) error {
		model.SetCursorMode(cursorMode)
		return nil
	}
}
