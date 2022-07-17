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

type Grammar interface {
}

type Model[T Grammar] struct {
	textinput  textinput.Model
	parser     *participle.Parser[T]
	parsedText *T
	tokens     []lexer.Token
	prompt     string
}

func New[T Grammar](parser *participle.Parser[T]) *Model[T] {
	textinput := textinput.New()
	textinput.Focus()
	return &Model[T]{parser: parser, textinput: textinput}
}

func (m *Model[T]) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model[T]) updateParsed() {
	expr, err := m.parser.ParseString("", m.Value(), participle.AllowTrailing(true))
	if err == nil {
		m.parsedText = expr
	} else {
		println(err.Error())
	}
	tokens, err := m.parser.Lex("", strings.NewReader(m.Value()))
	if err == nil {
		m.tokens = tokens
	} else {
		println(err.Error())
	}
}

func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	m.updateParsed()
	return cmd
}

func (m *Model[T]) View() string {
	lexer := lexers.Get("javascript")
	iter, err := lexer.Tokenise(nil, m.textinput.Value())
	style := styles.Get("swapoff")
	if err != nil {
		println(err)
	}
	return m.textinput.Prompt + m.inputFormatter(style, iter)
}

func (m *Model[T]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *Model[T]) Focused() bool {
	return m.textinput.Focused()
}

func (m *Model[T]) Parsed() *T {
	return m.parsedText
}

func (m *Model[T]) Value() string {
	return m.textinput.Value()
}

func (m *Model[T]) SetValue(value string) {
	m.textinput.SetValue(value)
}

func (m *Model[T]) Blur() {
	m.textinput.Blur()
}

func (m *Model[T]) Cursor() int {
	return m.textinput.Cursor()
}

func (m *Model[T]) SetCursor(cursor int) {
	m.textinput.SetCursor(cursor)
}

func (m *Model[T]) Prompt() string {
	return m.prompt
}

func (m *Model[T]) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m *Model[T]) ShouldSelectSuggestion(suggestion input.Suggestion[T]) bool {
	return suggestion.Text == m.CompletionText(m.Value())
}

func (m *Model[T]) currentToken(text string) (int, *lexer.Token) {
	cursor := m.Cursor()
	for i, token := range m.tokens {
		if cursor >= token.Pos.Offset && cursor <= token.Pos.Offset+len(token.String()) {
			return i, &token
		}
	}
	return -1, nil
}

func (m *Model[T]) CurrentToken() (int, *lexer.Token) {
	return m.currentToken(m.Value())
}

func (m *Model[T]) PreviousToken() (int, *lexer.Token) {
	index, _ := m.CurrentToken()
	if index <= 0 {
		return -1, nil
	}
	return index - 1, &m.tokens[index-1]
}

func (m *Model[T]) CompletionText(text string) string {
	_, token := m.currentToken(text)
	return token.String()
}

func (m *Model[T]) Tokens() []lexer.Token {
	return m.tokens
}

func (m *Model[T]) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[T]) tea.Cmd {
	return nil
}

func (m *Model[T]) OnSuggestionChanged(suggestion input.Suggestion[T]) {
	m.SetValue(suggestion.Text)
	m.SetCursor(len(suggestion.Text))
}

func (m *Model[T]) IsDelimiter(text string) bool {
	return false
}

func (m *Model[T]) OnSuggestionUnselected() {}

func (m *Model[T]) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	return false
}

func (m *Model[T]) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	return false
}

func (m *Model[T]) CurrentTokenBeforeCursor() string {
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
