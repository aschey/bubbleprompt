package formatter

import "github.com/charmbracelet/lipgloss"

type Formatters struct {
	Name              SuggestionText
	Description       SuggestionText
	ErrorText         lipgloss.Style
	SelectedIndicator lipgloss.Style
	Scrollbar         lipgloss.Style
	ScrollbarThumb    lipgloss.Style
	Suggestions       lipgloss.Style
}

var (
	DefaultNameForeground         = "243"
	DefaultNameBackground         = "7"
	DefaultSelectedNameForeground = "8"
	DefaultSelectedNameBackground = "14"
)

var (
	DefaultDescriptionForeground         = "255"
	DefaultDescriptionBackground         = "245"
	DefaultSelectedDescriptionForeground = "0"
	DefaultSelectedDescriptionBackground = "6"
	DefaultErrorTextBackground           = "1"
)

var (
	DefaultScrollbarColor      = "251"
	DefaultScrollbarThumbColor = "255"
)

var DefaultIndicatorForeground = "8"

func DefaultFormatters() Formatters {
	return Formatters{
		Name: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(DefaultSelectedNameForeground)).
				Background(lipgloss.Color(DefaultSelectedNameBackground)),
			Style: lipgloss.NewStyle().
				Foreground(lipgloss.Color(DefaultNameForeground)).
				Background(lipgloss.Color(DefaultNameBackground)),
		},
		Description: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(DefaultSelectedDescriptionForeground)).
				Background(lipgloss.Color(DefaultSelectedDescriptionBackground)),
			Style: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultDescriptionForeground)).
				Background(lipgloss.Color(DefaultDescriptionBackground)),
		},
		SelectedIndicator: lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(DefaultIndicatorForeground)),

		ErrorText: lipgloss.
			NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Background(lipgloss.Color(DefaultErrorTextBackground)),
		Scrollbar: lipgloss.
			NewStyle().
			Background(lipgloss.Color(DefaultScrollbarColor)),
		ScrollbarThumb: lipgloss.
			NewStyle().
			Background(lipgloss.Color(DefaultScrollbarThumbColor)),
	}
}

func (f Formatters) Minimal() Formatters {
	f.Name.Style = f.Name.Style.
		UnsetBackground().
		Foreground(lipgloss.Color(DefaultNameBackground))

	f.Name.SelectedStyle = f.Name.SelectedStyle.
		UnsetBackground().
		Foreground(lipgloss.Color(DefaultSelectedNameBackground))

	f.Description.Style = f.Description.Style.
		UnsetBackground().
		Foreground(lipgloss.Color(DefaultDescriptionBackground))

	f.Description.SelectedStyle = f.Description.SelectedStyle.
		UnsetBackground().
		Foreground(lipgloss.Color(DefaultSelectedDescriptionBackground))

	f.Scrollbar = f.Scrollbar.Background(lipgloss.Color("237"))
	f.ScrollbarThumb = f.ScrollbarThumb.Background(lipgloss.Color("240"))
	f.Suggestions = f.Suggestions.Border(lipgloss.RoundedBorder())

	return f
}
