package simpleinput

import (
	"github.com/aschey/bubbleprompt/editor/parser"
	"github.com/aschey/bubbleprompt/editor/parserinput"
	"github.com/charmbracelet/lipgloss"
)

type settings[T any] struct {
	delimiterRegex    string
	tokenRegex        string
	selectedTextStyle lipgloss.Style
	formatter         *parser.Formatter
	lexerOptions      []parserinput.Option[T]
}

type Option[T any] func(settings *settings[T]) error

func WithDelimiterRegex[T any](delimiterRegex string) Option[T] {
	return func(settings *settings[T]) error {
		settings.delimiterRegex = delimiterRegex
		return nil
	}
}

func WithTokenRegex[T any](tokenRegex string) Option[T] {
	return func(settings *settings[T]) error {
		settings.tokenRegex = tokenRegex
		return nil
	}
}

func WithSelectedTextStyle[T any](style lipgloss.Style) Option[T] {
	return func(settings *settings[T]) error {
		settings.selectedTextStyle = style
		return nil
	}
}

func WithFormatter[T any](formatter parser.Formatter) Option[T] {
	return func(settings *settings[T]) error {
		settings.formatter = &formatter
		return nil
	}
}

func WithLexerOptions[T any](options ...parserinput.Option[T]) Option[T] {
	return func(settings *settings[T]) error {
		settings.lexerOptions = options
		return nil
	}
}
