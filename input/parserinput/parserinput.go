package parserinput

import (
	"github.com/aschey/bubbleprompt/input/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type ParserModel[G any] struct {
	LexerModel
	parser     parser.Parser[G]
	parsedText *G
}

func NewParserModel[G any](parser parser.Parser[G], options ...Option) *ParserModel[G] {
	lexerModel := NewLexerModel(parser.Lexer(), options...)
	return &ParserModel[G]{parser: parser, LexerModel: *lexerModel}
}

func (m *ParserModel[G]) SetValue(value string) {
	m.LexerModel.SetValue(value)
	m.updateParsed()
}

func (m *ParserModel[G]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	cmd := m.LexerModel.OnUpdateStart(msg)
	m.updateParsed()
	return cmd
}

func (m *ParserModel[G]) Parsed() *G {
	return m.parsedText
}

func (m *ParserModel[G]) ParsedBeforeCursor() *G {
	expr, _ := m.parser.Parse(m.Value()[:m.Cursor()])
	return expr
}

func (m *ParserModel[G]) updateParsed() {
	expr, err := m.parser.Parse(m.Value())
	if err == nil {
		m.parsedText = expr
	} else {
		m.err = err
		return
	}
}
