package formatter

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
	_, hasNoBackground := t.Style.GetBackground().(lipgloss.NoColor)
	_, hasNoSelectedBackground := t.SelectedStyle.GetBackground().(lipgloss.NoColor)
	// Add left and right padding between each text section
	var leftPadding int
	if hasNoBackground && hasNoSelectedBackground {
		leftPadding = 0
	} else {
		leftPadding = 1
	}

	formattedText := style.
		Copy().
		PaddingLeft(leftPadding).
		PaddingRight(maxLen - runewidth.StringWidth(text) + 1).
		Render(text)
	return formattedText
}
