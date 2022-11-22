package lexerinput

import (
	"github.com/aschey/bubbleprompt/parser"
	"github.com/charmbracelet/bubbles/textinput"
)

type Option[T any] func(model *Model[T]) error

func WithDelimiterTokens[T any](tokens ...string) Option[T] {
	return func(model *Model[T]) error {
		model.delimiterTokens = tokens
		return nil
	}
}

func WithDelimiters[T any](delimiters ...string) Option[T] {
	return func(model *Model[T]) error {
		model.delimiters = delimiters
		return nil
	}
}

func WithFormatter[T any](formatter parser.Formatter) Option[T] {
	return func(model *Model[T]) error {
		model.formatter = formatter
		return nil
	}
}

func WithCursorMode[T any](cursorMode textinput.CursorMode) Option[T] {
	return func(model *Model[T]) error {
		model.SetCursorMode(cursorMode)
		return nil
	}
}
