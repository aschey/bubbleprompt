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

type Model[G any] struct {
	textinput  textinput.Model
	parser     *participle.Parser[G]
	parsedText *G
	tokens     []lexer.Token
	prompt     string
	err        error
}

func New[G any](parser *participle.Parser[G]) *Model[G] {
	textinput := textinput.New()
	textinput.Focus()
	return &Model[G]{parser: parser, textinput: textinput}
}

func (m *Model[G]) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model[G]) updateParsed() {
	tokens, lexErr := m.parser.Lex("", strings.NewReader(m.Value()))
	if lexErr == nil {
		m.tokens = tokens
	} else {
		m.err = lexErr
		return
	}

	expr, parseErr := m.parser.ParseString("", m.Value())
	if parseErr == nil {
		m.parsedText = expr
	} else {
		m.err = parseErr
		return
	}

	m.err = nil
}

func (m *Model[G]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	m.updateParsed()
	return cmd
}

func (m *Model[G]) Error() error {
	return m.err
}

func (m *Model[G]) View(viewMode input.ViewMode) string {
	lexer := lexers.Get("javascript")
	iter, err := lexer.Tokenise(nil, m.textinput.Value())
	style := styles.Get("swapoff")
	if err != nil {
		println(err)
	}
	return m.textinput.Prompt + m.inputFormatter(style, iter, viewMode)
}

func (m *Model[G]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *Model[G]) Focused() bool {
	return m.textinput.Focused()
}

func (m *Model[G]) Parsed() *G {
	return m.parsedText
}

func (m *Model[G]) ParsedBeforeCursor() *G {
	expr, _ := m.parser.ParseString("", m.Value()[:m.Cursor()])
	return expr
}

func (m *Model[G]) Value() string {
	return m.textinput.Value()
}

func (m *Model[G]) SetValue(value string) {
	m.textinput.SetValue(value)
	m.updateParsed()
}

func (m *Model[G]) Blur() {
	m.textinput.Blur()
}

func (m *Model[G]) Cursor() int {
	return m.textinput.Cursor()
}

func (m *Model[G]) SetCursor(cursor int) {
	m.textinput.SetCursor(cursor)
}

func (m *Model[G]) Prompt() string {
	return m.prompt
}

func (m *Model[G]) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m *Model[G]) ShouldSelectSuggestion(suggestion input.Suggestion[any]) bool {
	_, token := m.CurrentToken()
	tokenStr := token.Value
	return m.Cursor()-1 == token.Pos.Offset+len(tokenStr) && tokenStr == suggestion.Text
}

func (m *Model[G]) currentToken(text string, tokenPos int) (int, lexer.Token) {
	if len(m.tokens) == 0 {
		return -1, lexer.EOFToken(lexer.Position{})
	}
	if len(m.Value()) == 0 {
		return 0, m.tokens[0]
	}

	// Remove EOF token
	tokens := m.tokens[:len(m.tokens)-1]
	for i, token := range tokens {
		if i == len(tokens)-1 || (tokenPos >= token.Pos.Offset && tokenPos < token.Pos.Offset+len(token.Value)) {
			return i, token
		}
	}
	return -1, lexer.EOFToken(lexer.Position{})
}

func (m *Model[G]) CurrentToken() (int, lexer.Token) {
	return m.currentToken(m.Value(), m.Cursor()-1)
}

func (m *Model[G]) PreviousToken() (int, *lexer.Token) {
	index, _ := m.CurrentToken()
	if index <= 0 {
		return -1, nil
	}
	return index - 1, &m.tokens[index-1]
}

func (m *Model[G]) CompletionText(text string) string {
	_, token := m.currentToken(text, m.Cursor()-1)
	return token.String()
}

func (m *Model[G]) Tokens() []lexer.Token {
	return m.tokens
}

func (m *Model[G]) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[any]) tea.Cmd {
	return nil
}

func (m *Model[G]) OnSuggestionChanged(suggestion input.Suggestion[any]) {
	i, token := m.currentToken(m.Value(), m.Cursor())

	symbols := lexer.SymbolsByRune(m.parser.Lexer())
	tokenType := symbols[token.Type]
	if tokenType == "Punct" {
		if m.Cursor() < len(m.Value()) {
			token = m.tokens[i-1]
			if symbols[token.Type] == "Punct" {
				token = lexer.Token{Pos: lexer.Position{Offset: token.Pos.Offset + len(token.Value)}}
			}
		} else {
			token = m.tokens[i+1]
		}
	}

	start := token.Pos.Offset

	rest := start + len(token.Value)
	value := m.Value()
	newVal := value[:start] + suggestion.Text
	if rest < len(value) {
		newVal += value[start+len(token.Value):]
	}
	m.SetValue(newVal)
	m.SetCursor(start + len(suggestion.Text) - suggestion.CursorOffset)

}

func (m *Model[G]) OnSuggestionUnselected() {}

func (m *Model[G]) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	// Don't reset if no text because the completer won't run again
	return len(m.Value()) > 0
}

func (m *Model[G]) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	pos := m.Cursor()
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete:
		return pos < len(prevText)
	default:
		return true
	}
}

func (m *Model[G]) CurrentTokenBeforeCursor() string {
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

func (m *Model[G]) OnExecutorFinished() {
	// Clear out error once inpu text is reset
	m.err = nil
}
