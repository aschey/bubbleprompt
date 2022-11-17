package editor

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type SuggestionText struct {
	Style         lipgloss.Style
	SelectedStyle lipgloss.Style
}

func (t SuggestionText) Format(text string, maxLen int, selected bool) string {
	style := t.Style
	if selected {
		style = t.SelectedStyle
	}
	// Add left and right padding between each text section
	formattedText := style.
		Copy().
		PaddingLeft(1).
		PaddingRight(maxLen - runewidth.StringWidth(text) + 1).
		Render(text)
	return formattedText

}
