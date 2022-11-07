package input

import "github.com/charmbracelet/lipgloss"

type Formatters struct {
	Name               SuggestionText
	Description        SuggestionText
	DefaultPlaceholder lipgloss.Style
	ErrorText          lipgloss.Style
	Scrollbar          lipgloss.Style
	ScrollbarThumb     lipgloss.Style
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

var DefaultScrollbarColor = "251"
var DefaultScrollbarThumbColor = "255"

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
		DefaultPlaceholder: lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("6")),
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
