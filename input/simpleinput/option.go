package simpleinput

import (
	"regexp"

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

func WithDelimiterRegex(delimiterRegex regexp.Regexp) Option {
	return func(settings *settings) error {
		settings.delimiterRegex = delimiterRegex.String()
		return nil
	}
}

func WithTokenRegex(tokenRegex regexp.Regexp) Option {
	return func(settings *settings) error {
		settings.tokenRegex = tokenRegex.String()
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
