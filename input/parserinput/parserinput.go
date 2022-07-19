package parserinput

import (
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Token interface {
	SkipPrevious() bool
}

type Model[T Token, G any] struct {
	textinput  textinput.Model
	parser     *participle.Parser[G]
	parsedText *G
	tokens     []lexer.Token
	prompt     string
	err        error
}

func New[T Token, G any](parser *participle.Parser[G]) *Model[T, G] {
	textinput := textinput.New()
	textinput.Focus()
	return &Model[T, G]{parser: parser, textinput: textinput}
}

func (m *Model[T, G]) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model[T, G]) updateParsed() {
	tokens, lexErr := m.parser.Lex("", strings.NewReader(m.Value()))
	if lexErr == nil {
		m.tokens = tokens
	} else {
		m.err = lexErr
		return
	}

	expr, parseErr := m.parser.ParseString("", m.Value(), participle.AllowTrailing(true))
	if parseErr == nil {
		m.parsedText = expr
	} else {
		m.err = parseErr
		return
	}

	m.err = nil
}

func (m *Model[T, G]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	m.updateParsed()
	return cmd
}

func (m *Model[T, G]) Error() error {
	return m.err
}

func (m *Model[T, G]) View(viewMode input.ViewMode) string {
	lexer := lexers.Get("javascript")
	iter, err := lexer.Tokenise(nil, m.textinput.Value())
	style := styles.Get("swapoff")
	if err != nil {
		println(err)
	}
	return m.textinput.Prompt + m.inputFormatter(style, iter, viewMode)
}

func (m *Model[T, G]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *Model[T, G]) Focused() bool {
	return m.textinput.Focused()
}

func (m *Model[T, G]) Parsed() *G {
	return m.parsedText
}

func (m *Model[T, G]) Value() string {
	return m.textinput.Value()
}

func (m *Model[T, G]) SetValue(value string) {
	m.textinput.SetValue(value)
	m.updateParsed()
}

func (m *Model[T, G]) Blur() {
	m.textinput.Blur()
}

func (m *Model[T, G]) Cursor() int {
	return m.textinput.Cursor()
}

func (m *Model[T, G]) SetCursor(cursor int) {
	m.textinput.SetCursor(cursor)
}

func (m *Model[T, G]) Prompt() string {
	return m.prompt
}

func (m *Model[T, G]) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m *Model[T, G]) ShouldSelectSuggestion(suggestion input.Suggestion[T]) bool {
	return suggestion.Text == m.CompletionText(m.Value())
}

func (m *Model[T, G]) currentToken(text string) (int, *lexer.Token) {
	cursor := m.Cursor()
	for i, token := range m.tokens {
		if cursor >= token.Pos.Offset && cursor <= token.Pos.Offset+len(token.String()) {
			return i, &token
		}
	}
	return -1, nil
}

func (m *Model[T, G]) CurrentToken() (int, *lexer.Token) {
	return m.currentToken(m.Value())
}

func (m *Model[T, G]) PreviousToken() (int, *lexer.Token) {
	index, _ := m.CurrentToken()
	if index <= 0 {
		return -1, nil
	}
	return index - 1, &m.tokens[index-1]
}

func (m *Model[T, G]) CompletionText(text string) string {
	_, token := m.currentToken(text)
	return token.String()
}

func (m *Model[T, G]) Tokens() []lexer.Token {
	return m.tokens
}

func (m *Model[T, G]) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[T]) tea.Cmd {
	return nil
}

func (m *Model[T, G]) OnSuggestionChanged(suggestion input.Suggestion[T]) {
	_, token := m.CurrentToken()
	if token == nil {
		return
	}
	start := token.Pos.Offset

	if !token.EOF() && suggestion.Metadata.SkipPrevious() {
		start += len(token.String())
	}

	m.SetValue(m.Value()[:start] + suggestion.Text)
	m.SetCursor(start + len(suggestion.Text) - suggestion.CursorOffset)

}

func (m *Model[T, G]) OnSuggestionUnselected() {}

func (m *Model[T, G]) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	return false
}

func (m *Model[T, G]) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	return false
}

func (m *Model[T, G]) CurrentTokenBeforeCursor() string {
	if m.parsedText == nil {
		return ""
	}
	_, token := m.CurrentToken()
	start := token.Pos.Offset
	cursor := m.Cursor()
	if start > cursor {
		return ""
	}
	val := m.Value()[start:cursor]
	return val
}

func (m *Model[T, G]) OnExecutorFinished() {
	// Clear out error once inpu text is reset
	m.err = nil
}
