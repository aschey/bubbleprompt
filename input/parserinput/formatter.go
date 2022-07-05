package parserinput

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/charmbracelet/lipgloss"
)

func (m Model[T]) inputFormatter(theme *chroma.Style, iter chroma.Iterator) string {
	theme = clearBackground(theme)
	formatted := ""
	pos := 0
	cursor := m.textinput.Cursor()
	blink := m.textinput.Blink()
	for token := iter(); token != chroma.EOF; token = iter() {
		entry := theme.Get(token.Type)
		style := lipgloss.NewStyle()
		if !entry.IsZero() {
			if entry.Bold == chroma.Yes {
				style = style.Bold(true)
			}
			if entry.Underline == chroma.Yes {
				style = style.Underline(true)
			}
			if entry.Italic == chroma.Yes {
				style = style.Italic(true)
			}
			if entry.Colour.IsSet() {
				style = style.Foreground(lipgloss.Color(entry.Colour.String()))
			}
			if entry.Background.IsSet() {
				style = style.Background(lipgloss.Color(entry.Background.String()))
			}
		}

		if cursor >= pos && cursor < pos+len(token.Value) {
			localCursor := cursor - pos
			formatted += style.Render(token.Value[:localCursor])
			cursorStyle := style.Copy().Reverse(!blink)
			formatted += cursorStyle.Render(string(token.Value[localCursor]))
			if localCursor < len(token.Value)-1 {
				formatted += style.Render(token.Value[localCursor+1:])
			}
		} else {
			formatted += style.Render(token.Value)
		}

		pos += len(token.Value)
	}
	if cursor >= pos {
		formatted += lipgloss.NewStyle().Reverse(!blink).Render(" ")
	}
	return formatted
}

func clearBackground(style *chroma.Style) *chroma.Style {
	builder := style.Builder()
	bg := builder.Get(chroma.Background)
	bg.Background = 0
	bg.NoInherit = true
	builder.AddEntry(chroma.Background, bg)
	style, _ = builder.Build()
	return style
}
