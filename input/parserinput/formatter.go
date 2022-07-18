package parserinput

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/lipgloss"
)

func (m Model[T, G]) inputFormatter(theme *chroma.Style, iter chroma.Iterator) string {
	theme = clearBackground(theme)
	viewBuilder := input.NewViewBuilder(m.Cursor(), lipgloss.NewStyle(), " ", m.textinput.Blink())
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
		viewBuilder.Render(strings.TrimRight(token.Value, "\n"), viewBuilder.ViewLen(), style)
	}

	return viewBuilder.GetView()
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
