package parserinput

import (
	"github.com/aschey/bubbleprompt/input/lexerinput"
	"github.com/aschey/bubbleprompt/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type Model[T any, G any] struct {
	lexerinput.Model[T]
	parser     parser.Parser[G]
	parsedText *G
	err        error
}

func NewModel[T any, G any](parser parser.Parser[G], options ...lexerinput.Option[T]) *Model[T, G] {
	lexerModel := lexerinput.NewModel(parser.Lexer(), options...)
	return &Model[T, G]{parser: parser, Model: *lexerModel}
}

func (m *Model[T, G]) SetValue(value string) {
	m.Model.SetValue(value)
	m.updateParsed()
}

func (m *Model[T, G]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	cmd := m.Model.OnUpdateStart(msg)
	m.updateParsed()
	return cmd
}

func (m Model[T, G]) Parsed() *G {
	return m.parsedText
}

func (m Model[T, G]) ParsedBeforeCursor() (*G, error) {
	return m.parser.Parse(string(m.Runes()[:m.CursorIndex()]))
}

func (m *Model[T, G]) updateParsed() {
	expr, err := m.parser.Parse(m.Value())
	if err == nil {
		m.parsedText = expr
	} else {
		m.err = err
		return
	}
}

func (m Model[T, G]) Error() error {
	if m.err != nil {
		return m.err
	}
	return m.Model.Error()
}
