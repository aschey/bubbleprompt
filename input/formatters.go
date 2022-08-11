package input

import "github.com/charmbracelet/lipgloss"

type Formatters struct {
	Name               SuggestionText
	Description        SuggestionText
	DefaultPlaceholder Text
	ErrorText          Text
}

const DefaultNameBackground = "12"
const DefaultDescriptionBackground = "13"
const DefaultSelectedForeground = "8"
const DefaultErrorTextBackground = "1"

func DefaultFormatters() Formatters {
	return Formatters{
		Name: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultSelectedForeground)).
				Background(lipgloss.Color(DefaultNameBackground)),
			Style: lipgloss.NewStyle().Background(lipgloss.Color(DefaultNameBackground)),
		},
		Description: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultSelectedForeground)).
				Background(lipgloss.Color(DefaultDescriptionBackground)),
			Style: lipgloss.
				NewStyle().
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
