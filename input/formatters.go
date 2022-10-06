package input

import "github.com/charmbracelet/lipgloss"

type Formatters struct {
	Name               SuggestionText
	Description        SuggestionText
	DefaultPlaceholder Text
	ErrorText          Text
}

var DefaultNameForeground = "255"
var DefaultNameBackground = "6"
var DefaultSelectedNameForeground = "240"
var DefaultSelectedNameBackground = "14"

var DefaultDescriptionForeground = "255"
var DefaultDescriptionBackground = "245"
var DefaultSelectedDescriptionForeground = "240"
var DefaultSelectedDescriptionBackground = "249"
var DefaultErrorTextBackground = "1"

func DefaultFormatters() Formatters {
	return Formatters{
		Name: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultSelectedNameForeground)).
				Background(lipgloss.Color(DefaultSelectedNameBackground)),
			Style: lipgloss.NewStyle().
				Foreground(lipgloss.Color(DefaultNameForeground)).
				Background(lipgloss.Color(DefaultNameBackground)),
		},
		Description: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultSelectedDescriptionForeground)).
				Background(lipgloss.Color(DefaultSelectedDescriptionBackground)),
			Style: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultDescriptionForeground)).
				Background(lipgloss.Color(DefaultDescriptionBackground)),
		},
		DefaultPlaceholder: Text{
			Style: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color("6")),
		},
		ErrorText: Text{
			Style: lipgloss.
				NewStyle().
				PaddingLeft(1).
				PaddingRight(1).
				Background(lipgloss.Color(DefaultErrorTextBackground)),
		},
	}
}
