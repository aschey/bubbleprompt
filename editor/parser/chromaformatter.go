package parser

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/aschey/bubbleprompt/editor"
	"github.com/charmbracelet/lipgloss"
)

type ChromaFormatter struct {
	style *chroma.Style
	lexer chroma.Lexer
}

func NewChromaFormatter(style *chroma.Style, lexer chroma.Lexer) *ChromaFormatter {
	return &ChromaFormatter{style: style, lexer: lexer}
}

func (c *ChromaFormatter) Lex(input string, _ *editor.Token) ([]FormatterToken, error) {
	theme := clearBackground(c.style)
	iter, err := c.lexer.Tokenise(nil, input)
	if err != nil {
		return nil, err
	}
	tokens := []FormatterToken{}
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
		tokens = append(tokens, FormatterToken{Value: token.Value, Style: style})
	}

	return tokens, nil
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
