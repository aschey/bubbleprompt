package prompt

import "github.com/charmbracelet/lipgloss"

type Formatter func(name string, columnWidth int, selected bool) string

type Text struct {
	ForegroundColor string
	BackgroundColor string
	Formatter       func(text string) string
}

func (p Text) format(text string) string {
	if p.Formatter == nil {
		return lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(p.ForegroundColor)).
			Background(lipgloss.Color(p.BackgroundColor)).
			Render(text)
	}
	return p.Formatter(text)
}

type SuggestionText struct {
	ForegroundColor         string
	SelectedForegroundColor string
	BackgroundColor         string
	SelectedBackgroundColor string
	Formatter               Formatter
}

func (t SuggestionText) format(text string, maxLen int, selected bool) string {
	defaultStyle := lipgloss.
		NewStyle().
		PaddingLeft(1)

	foregroundColor := t.ForegroundColor
	backgroundColor := t.BackgroundColor
	if selected {
		foregroundColor = t.SelectedForegroundColor
		backgroundColor = t.SelectedBackgroundColor
	}
	var formattedText string
	if t.Formatter == nil {
		formattedText = defaultStyle.
			Copy().
			Foreground(lipgloss.Color(foregroundColor)).
			Background(lipgloss.Color(backgroundColor)).
			PaddingRight(maxLen - len(text) + 1).
			Render(text)
	} else {
		formattedText = t.Formatter(text, maxLen, selected)
	}

	return formattedText
}
