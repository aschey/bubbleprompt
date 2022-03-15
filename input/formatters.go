package input

import "github.com/charmbracelet/lipgloss"

type Formatters struct {
	Name               SuggestionText
	Description        SuggestionText
	DefaultPlaceholder Text
	SelectedSuggestion lipgloss.Style
}
