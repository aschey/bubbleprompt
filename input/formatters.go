package input

import "github.com/charmbracelet/lipgloss"

type Formatters struct {
	Name               SuggestionText
	Description        SuggestionText
	DefaultPlaceholder Text
	ErrorText          Text
}

const DefaultNameBackground = "14"
const DefaultDescriptionBackground = "37"

func DefaultFormatters() Formatters {
	return Formatters{
		Name: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color("240")).
				Background(lipgloss.Color(DefaultNameBackground)),
			Style: lipgloss.NewStyle().Background(lipgloss.Color("14")),
		},
		Description: SuggestionText{
			SelectedStyle: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color("240")).
				Background(lipgloss.Color(DefaultDescriptionBackground)),
			Style: lipgloss.NewStyle().Background(lipgloss.Color("37")),
		},
		DefaultPlaceholder: Text{
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		},
		ErrorText: Text{
			Style: lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Background(lipgloss.Color("#ff0000")),
		},
	}
}
