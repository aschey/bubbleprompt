package simpleinput

import (
	"github.com/aschey/bubbleprompt/input/lexerinput"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/charmbracelet/lipgloss"
)

type settings[T any] struct {
	delimiterRegex    string
	tokenRegex        string
	selectedTextStyle lipgloss.Style
	formatter         *parser.Formatter
	lexerOptions      []lexerinput.Option[T]
}

type Option[T any] func(settings *settings[T])

func WithDelimiterRegex[T any](delimiterRegex string) Option[T] {
	return func(settings *settings[T]) {
		settings.delimiterRegex = delimiterRegex
	}
}

func WithTokenRegex[T any](tokenRegex string) Option[T] {
	return func(settings *settings[T]) {
		settings.tokenRegex = tokenRegex
	}
}

func WithSelectedTextStyle[T any](style lipgloss.Style) Option[T] {
	return func(settings *settings[T]) {
		settings.selectedTextStyle = style
	}
}

func WithFormatter[T any](formatter parser.Formatter) Option[T] {
	return func(settings *settings[T]) {
		settings.formatter = &formatter
	}
}

func WithLexerOptions[T any](options ...lexerinput.Option[T]) Option[T] {
	return func(settings *settings[T]) {
		settings.lexerOptions = options
	}
}
