package simpleinput

import (
	"github.com/aschey/bubbleprompt/input/lexerinput"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"
)

type settings[T any] struct {
	delimiterRegex    string
	tokenRegex        string
	selectedTextStyle lipgloss.Style
	formatterFunc     func(parser.Lexer) parser.Formatter
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

func WithFormatter[T any](formatterFunc func(lexer parser.Lexer) parser.Formatter) Option[T] {
	return func(settings *settings[T]) {
		settings.formatterFunc = formatterFunc
	}
}

func WithCursorMode[T any](cursorMode cursor.Mode) Option[T] {
	return func(settings *settings[T]) {
		settings.lexerOptions = append(
			settings.lexerOptions,
			lexerinput.WithCursorMode[T](cursorMode),
		)
	}
}

func WithPrompt[T any](prompt string) Option[T] {
	return func(settings *settings[T]) {
		settings.lexerOptions = append(
			settings.lexerOptions,
			lexerinput.WithPrompt[T](prompt),
		)
	}
}
