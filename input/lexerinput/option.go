package lexerinput

import (
	"github.com/aschey/bubbleprompt/parser"
	"github.com/charmbracelet/bubbles/cursor"
)

type Option[T any] func(model *Model[T])

func WithDelimiterTokens[T any](tokens ...string) Option[T] {
	return func(model *Model[T]) {
		model.delimiterTokens = tokens
	}
}

func WithDelimiters[T any](delimiters ...string) Option[T] {
	return func(model *Model[T]) {
		model.delimiters = delimiters
	}
}

func WithFormatter[T any](formatter parser.Formatter) Option[T] {
	return func(model *Model[T]) {
		model.formatter = formatter
	}
}

func WithCursorMode[T any](cursorMode cursor.Mode) Option[T] {
	return func(model *Model[T]) {
		model.SetCursorMode(cursorMode)
	}
}
