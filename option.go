package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/lipgloss"
)

type Option func(model *Model) error

func WithNameStyle(style lipgloss.Style) Option {
	return func(model *Model) error {
		model.Formatters.Name.Style = style
		return nil
	}
}

func WithSelectedNameStyle(style lipgloss.Style) Option {
	return func(model *Model) error {
		model.Formatters.Name.SelectedStyle = style
		return nil
	}
}

func WithNameFormatter(nameFormatter input.Formatter) Option {
	return func(model *Model) error {
		model.Formatters.Name.Formatter = nameFormatter
		return nil
	}
}

func WithDescriptionStyle(style lipgloss.Style) Option {
	return func(model *Model) error {
		model.Formatters.Description.Style = style
		return nil
	}
}

func WithSelectedDescriptionStyle(style lipgloss.Style) Option {
	return func(model *Model) error {
		model.Formatters.Description.SelectedStyle = style
		return nil
	}
}

func WithDescriptionFormatter(descriptionFormatter input.Formatter) Option {
	return func(model *Model) error {
		model.Formatters.Description.Formatter = descriptionFormatter
		return nil
	}
}

func WithDefaultPlaceholderStyle(style lipgloss.Style) Option {
	return func(model *Model) error {
		model.Formatters.DefaultPlaceholder.Style = style
		return nil
	}
}

func WithDefaultPlaceholderFormatter(formatter func(text string) string) Option {
	return func(model *Model) error {
		model.Formatters.DefaultPlaceholder.Formatter = formatter
		return nil
	}
}

func WithSelectedSuggestionStyle(style lipgloss.Style) Option {
	return func(model *Model) error {
		model.Formatters.SelectedSuggestion = style
		return nil
	}
}
