package simpleinput

import (
	"github.com/aschey/bubbleprompt/input/parser"
	"github.com/charmbracelet/lipgloss"
)

type settings struct {
	delimiterRegex    string
	tokenRegex        string
	selectedTextStyle lipgloss.Style
	formatter         *parser.Formatter
}

type Option func(settings *settings) error

func WithDelimiterRegex(delimiterRegex string) Option {
	return func(settings *settings) error {
		settings.delimiterRegex = delimiterRegex
		return nil
	}
}

func WithTokenRegex(tokenRegex string) Option {
	return func(settings *settings) error {
		settings.tokenRegex = tokenRegex
		return nil
	}
}

func WithSelectedTextStyle(style lipgloss.Style) Option {
	return func(settings *settings) error {
		settings.selectedTextStyle = style
		return nil
	}
}

func WithFormatter(formatter parser.Formatter) Option {
	return func(settings *settings) error {
		settings.formatter = &formatter
		return nil
	}
}
