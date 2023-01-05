package formatter

import "github.com/charmbracelet/lipgloss"

type Formatters struct {
	Name              SuggestionText
	Description       SuggestionText
	ErrorText         lipgloss.Style
	SelectedIndicator lipgloss.Style
	Scrollbar         lipgloss.Style
	ScrollbarThumb    lipgloss.Style
}

var DefaultNameForeground = "243"
var DefaultNameBackground = "7"
var DefaultSelectedNameForeground = "8"
var DefaultSelectedNameBackground = "14"

var DefaultDescriptionForeground = "255"
var DefaultDescriptionBackground = "245"
var DefaultSelectedDescriptionForeground = "0"
var DefaultSelectedDescriptionBackground = "6"
var DefaultErrorTextBackground = "1"

var DefaultScrollbarColor = "251"
var DefaultScrollbarThumbColor = "255"

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
	f.Name.Style.UnsetBackground()
	f.Name.SelectedStyle.UnsetBackground()
	f.Name.Style.Foreground(lipgloss.Color(DefaultNameBackground))
	f.Name.SelectedStyle.Foreground(lipgloss.Color(DefaultSelectedNameBackground))

	f.Description.Style.UnsetBackground()
	f.Description.SelectedStyle.UnsetBackground()
	f.Description.Style.Foreground(lipgloss.Color(DefaultDescriptionBackground))
	f.Description.SelectedStyle.Foreground(lipgloss.Color(DefaultSelectedDescriptionBackground))

	f.Scrollbar.Background(lipgloss.Color("237"))
	f.ScrollbarThumb.Background(lipgloss.Color("240"))

	return f
}