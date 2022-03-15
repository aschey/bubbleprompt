package input

import "github.com/charmbracelet/lipgloss"

type Formatter func(name string, columnWidth int, selected bool) string

type Text struct {
	Style     lipgloss.Style
	Formatter func(text string) string
}

func (p Text) Format(text string) string {
	if p.Formatter == nil {
		return p.Style.Render(text)
	}
	return p.Formatter(text)
}

type SuggestionText struct {
	Style         lipgloss.Style
	SelectedStyle lipgloss.Style
	Formatter     Formatter
}

func (t SuggestionText) Format(text string, maxLen int, selected bool) string {

	if t.Formatter == nil {
		style := t.Style
		if selected {
			style = t.SelectedStyle
		}
		// Add left and right padding between each text section
		formattedText := style.
			Copy().
			PaddingLeft(1).
			PaddingRight(maxLen - len(text) + 1).
			Render(text)
		return formattedText
	} else {
		formattedText := t.Formatter(text, maxLen, selected)
		return formattedText
	}

}
