package parserinput

import (
	"math"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

const defaultWhitespace = math.MinInt

type LexerModel struct {
	textinput       textinput.Model
	lexer           lexer.Definition
	tokens          []lexer.Token
	delimiterTokens []string
	delimiters      []string
	prompt          string
	err             error
}

func NewLexerModel(def lexer.Definition, options ...Option) *LexerModel {
	textinput := textinput.New()
	textinput.Focus()
	model := &LexerModel{lexer: def, textinput: textinput, tokens: []lexer.Token{}}
	for _, option := range options {
		if err := option(model); err != nil {
			panic(err)
		}
	}
	return model
}

func (m *LexerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LexerModel) updateTokens() {
	lex, err := m.lexer.Lex("", strings.NewReader(m.Value()))
	if err != nil {
		m.err = err
		return
	}
	tokens, err := lexer.ConsumeAll(lex)
	fullTokens := []lexer.Token{}
	if err == nil {
		for i, token := range tokens {
			if i > 0 {
				prevEnd := tokens[i-1].Pos.Offset + len(tokens[i-1].Value)
				if prevEnd < token.Pos.Offset {
					// This part of the input was ignored by the lexer
					// so insert a dummy token to account for it
					fullTokens = append(fullTokens, lexer.Token{
						Value: strings.Repeat(" ", token.Pos.Offset-prevEnd),
						Type:  defaultWhitespace,
						Pos: lexer.Position{
							Offset: prevEnd,
						}})
				}
			}
			fullTokens = append(fullTokens, token)
		}
		m.tokens = fullTokens
	} else {
		m.err = err
		return
	}

	m.err = nil
}

func (m *LexerModel) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	m.updateTokens()
	return cmd
}

func (m *LexerModel) Error() error {
	return m.err
}

func (m *LexerModel) View(viewMode input.ViewMode) string {
	lexer := lexers.Get("javascript")
	iter, err := lexer.Tokenise(nil, m.textinput.Value())
	style := styles.Get("swapoff")
	if err != nil {
		println(err)
	}
	return m.textinput.Prompt + m.inputFormatter(style, iter, viewMode)
}

func (m *LexerModel) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *LexerModel) Focused() bool {
	return m.textinput.Focused()
}

func (m *LexerModel) Value() string {
	return m.textinput.Value()
}

func (m *LexerModel) SetValue(value string) {
	m.textinput.SetValue(value)
	m.updateTokens()
}

func (m *LexerModel) Blur() {
	m.textinput.Blur()
}

func (m *LexerModel) Cursor() int {
	return m.textinput.Cursor()
}

func (m *LexerModel) SetCursor(cursor int) {
	m.textinput.SetCursor(cursor)
}

func (m *LexerModel) Prompt() string {
	return m.prompt
}

func (m *LexerModel) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m *LexerModel) ShouldSelectSuggestion(suggestion input.Suggestion[any]) bool {
	_, token := m.CurrentToken()
	tokenStr := token.Value
	return m.Cursor()-1 == token.Pos.Offset+len(tokenStr) && tokenStr == suggestion.Text
}

func (m *LexerModel) currentToken(text string, tokenPos int) (int, lexer.Token) {
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

func (m *LexerModel) CurrentToken() (int, lexer.Token) {
	return m.currentToken(m.Value(), m.Cursor()-1)
}

func (m *LexerModel) PreviousToken() (int, *lexer.Token) {
	index, _ := m.CurrentToken()
	if index <= 0 {
		return -1, nil
	}
	return index - 1, &m.tokens[index-1]
}

func (m *LexerModel) CompletionText(text string) string {
	_, token := m.currentToken(text, m.Cursor()-1)
	return token.String()
}

func (m *LexerModel) Tokens() []lexer.Token {
	return m.tokens
}

func (m *LexerModel) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[any]) tea.Cmd {
	return nil
}

func (m *LexerModel) isDelimiterToken(token lexer.Token) bool {
	symbols := lexer.SymbolsByRune(m.lexer)
	symbol := symbols[token.Type]
	// Dummy whitespace tokens won't be registered with the lexer so check it separately
	return slices.Contains(m.delimiters, token.Value) || slices.Contains(m.delimiterTokens, symbol) || token.Type == defaultWhitespace
}

func (m *LexerModel) OnSuggestionChanged(suggestion input.Suggestion[any]) {
	i, token := m.currentToken(m.Value(), m.Cursor())

	if m.isDelimiterToken(token) {
		if m.Cursor() < len(m.Value()) {
			token = m.tokens[i-1]
			if m.isDelimiterToken(token) {
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

func (m *LexerModel) OnSuggestionUnselected() {}

func (m *LexerModel) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	// Don't reset if no text because the completer won't run again
	return len(m.Value()) > 0
}

func (m *LexerModel) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	pos := m.Cursor()
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete:
		return pos < len(prevText)
	default:
		return true
	}
}

func (m *LexerModel) CompletableTokenBeforeCursor() string {
	_, token := m.CurrentToken()
	if m.isDelimiterToken(token) {
		// Don't filter suggestions on delimiters
		return ""
	}
	return m.currentTokenBeforeCursor(token)
}

func (m *LexerModel) CurrentTokenBeforeCursor() string {
	_, token := m.CurrentToken()

	return m.currentTokenBeforeCursor(token)
}

func (m *LexerModel) currentTokenBeforeCursor(token lexer.Token) string {
	start := token.Pos.Offset
	cursor := m.Cursor()
	if start > cursor {
		return ""
	}
	val := m.Value()[start:cursor]
	return val
}

func (m *LexerModel) OnExecutorFinished() {
	// Clear out error once input text is reset
	m.err = nil
}
