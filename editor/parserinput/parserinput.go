package parserinput

import (
	"github.com/aschey/bubbleprompt/editor/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type ParserModel[T any, G any] struct {
	LexerModel[T]
	parser     parser.Parser[G]
	parsedText *G
}

func NewParserModel[T any, G any](parser parser.Parser[G], options ...Option[T]) *ParserModel[T, G] {
	lexerModel := NewLexerModel(parser.Lexer(), options...)
	return &ParserModel[T, G]{parser: parser, LexerModel: *lexerModel}
}

func (m *ParserModel[T, G]) SetValue(value string) {
	m.LexerModel.SetValue(value)
	m.updateParsed()
}

func (m *ParserModel[T, G]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	cmd := m.LexerModel.OnUpdateStart(msg)
	m.updateParsed()
	return cmd
}

func (m *ParserModel[T, G]) Parsed() *G {
	return m.parsedText
}

func (m *ParserModel[T, G]) ParsedBeforeCursor() (*G, error) {
	return m.parser.Parse(string(m.Runes()[:m.CursorIndex()]))
}

func (m *ParserModel[T, G]) updateParsed() {
	expr, err := m.parser.Parse(m.Value())
	if err == nil {
		m.parsedText = expr
	} else {
		m.err = err
		return
	}
}
