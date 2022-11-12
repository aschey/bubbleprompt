package parserinput

import (
	"github.com/aschey/bubbleprompt/input/parser"
)

type Option[T any] func(model *LexerModel[T]) error

func WithDelimiterTokens[T any](tokens ...string) Option[T] {
	return func(model *LexerModel[T]) error {
		model.delimiterTokens = tokens
		return nil
	}
}

func WithDelimiters[T any](delimiters ...string) Option[T] {
	return func(model *LexerModel[T]) error {
		model.delimiters = delimiters
		return nil
	}
}

func WithFormatter[T any](formatter parser.Formatter) Option[T] {
	return func(model *LexerModel[T]) error {
		model.formatter = formatter
		return nil
	}
}
