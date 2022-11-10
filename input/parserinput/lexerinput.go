package parserinput

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parser"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"
)

type LexerModel struct {
	textinput        textinput.Model
	lexer            parser.Lexer
	formatter        parser.Formatter
	selectedToken    *parser.Token
	tokens           []parser.Token
	formatterTokens  []parser.FormatterToken
	delimiterTokens  []string
	delimiters       []string
	whitespaceTokens map[int]bool
	prompt           string
	err              error
}

func NewLexerModel(lexer parser.Lexer, options ...Option) *LexerModel {
	textinput := textinput.New()
	textinput.Focus()
	model := &LexerModel{
		lexer:            lexer,
		textinput:        textinput,
		prompt:           "> ",
		tokens:           []parser.Token{},
		formatterTokens:  []parser.FormatterToken{},
		whitespaceTokens: make(map[int]bool),
	}
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

func (m *LexerModel) createWhitespaceToken(start int, end int, index int) parser.Token {
	token := parser.Token{
		Value: m.Value()[start:end],
		Start: start,
		Index: index,
	}
	m.whitespaceTokens[start] = true
	return token
}

func (m *LexerModel) updateTokens() error {
	tokens, err := m.lexer.Lex(m.Value())
	if err != nil {
		return err
	}

	if m.formatter != nil {
		m.formatterTokens, err = m.formatter.Lex(m.Value(), m.selectedToken)
		if err != nil {
			return err
		}
	}

	fullTokens := []parser.Token{}
	m.whitespaceTokens = make(map[int]bool)
	last := 0
	index := 0
	for i, token := range tokens {
		if i > 0 {
			prevEnd := tokens[i-1].End()
			if prevEnd < token.Start {
				// This part of the input was ignored by the lexer
				// so insert a dummy token to account for it
				fullTokens = append(fullTokens, m.createWhitespaceToken(prevEnd, token.Start, index))
				index++
			}
		}
		token.Index = index
		index++
		fullTokens = append(fullTokens, token)
		last = token.End()
	}

	// Check for trailing whitespace
	if m.Cursor() > last {
		fullTokens = append(fullTokens, m.createWhitespaceToken(last, m.Cursor(), index))
	}
	m.tokens = fullTokens

	return nil
}

func (m *LexerModel) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	if msg, ok := msg.(tea.KeyMsg); ok {
		err := m.updateTokens()
		// Don't reset error on submit yet because we need to pass it to the view
		if msg.Type != tea.KeyEnter {
			m.err = err
		}
	}

	return cmd
}

func (m *LexerModel) Error() error {
	return m.err
}

func (m *LexerModel) unstyledView(text string, showCursor bool) string {
	viewBuilder := input.NewViewBuilder(m.Cursor(), lipgloss.NewStyle(), " ", showCursor)
	viewBuilder.Render(text, 0, lipgloss.NewStyle())
	return m.prompt + viewBuilder.View()
}

func (m *LexerModel) styledView(formatterTokens []parser.FormatterToken, showCursor bool) string {
	viewBuilder := input.NewViewBuilder(m.Cursor(), lipgloss.NewStyle(), " ", showCursor)
	for _, token := range formatterTokens {
		viewBuilder.Render(strings.TrimRight(token.Value, "\n"), viewBuilder.ViewLen(), token.Style)
	}

	return m.prompt + viewBuilder.View()
}

func (m *LexerModel) View(viewMode input.ViewMode) string {
	showCursor := !m.textinput.Blink()
	if viewMode == input.Static {
		showCursor = false
	}
	if m.formatter == nil {
		return m.unstyledView(m.Value(), showCursor)
	}

	return m.styledView(m.formatterTokens, showCursor)
}

func (m *LexerModel) FormatText(text string) string {
	if m.formatter == nil {
		return m.unstyledView(text, false)
	}
	formatterTokens, _ := m.formatter.Lex(text, nil)
	return m.styledView(formatterTokens, false)
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

func (m *LexerModel) ResetValue() {
	m.textinput.SetValue("")
	_ = m.updateTokens()
}

func (m *LexerModel) SetValue(value string) {
	m.textinput.SetValue(value)
	m.err = m.updateTokens()
}

func (m *LexerModel) setSelectedToken(token *parser.Token) {
	m.selectedToken = token
	_ = m.updateTokens()
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
	token := m.CurrentToken()
	tokenStr := token.Value
	return m.Cursor() == token.End() && tokenStr == suggestion.Text
}

func (m *LexerModel) currentToken(text string, tokenPos int) parser.Token {
	if len(m.tokens) == 0 {
		return parser.Token{Index: -1}
	}
	if len(m.Value()) == 0 {
		return m.tokens[0]
	}

	for i, token := range m.tokens {
		if i == len(m.tokens)-1 || (tokenPos >= token.Start && tokenPos < token.End()) {
			return token
		}
	}
	return parser.Token{Index: -1}
}

func (m *LexerModel) CurrentToken() parser.Token {
	return m.currentToken(m.Value(), m.Cursor()-1)
}

func (m *LexerModel) FindLast(filter func(token parser.Token, symbol string) bool) *parser.Token {
	currentToken := m.CurrentToken()
	for i := currentToken.Index; i >= 0; i-- {
		token := m.tokens[i]

		if filter(token, token.Type) {
			return &token
		}
	}

	return nil
}

func (m *LexerModel) PreviousToken() *parser.Token {
	currentToken := m.CurrentToken()
	if currentToken.Index <= 0 {
		return nil
	}
	return &m.tokens[currentToken.Index-1]
}

func (m *LexerModel) CompletionText(text string) string {
	token := m.currentToken(text, m.Cursor()-1)
	return token.Value
}

func (m *LexerModel) Tokens() []parser.Token {
	return m.tokens
}

func (m *LexerModel) TokenValues() []string {
	tokens := []string{}
	for _, token := range m.tokens {
		tokens = append(tokens, token.Value)
	}
	return tokens
}

func (m *LexerModel) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[any], isSelected bool) tea.Cmd {
	return nil
}

func (m *LexerModel) IsDelimiterToken(token parser.Token) bool {
	// Dummy whitespace tokens won't be registered with the lexer so check them separately
	return slices.Contains(m.delimiters, token.Value) || slices.Contains(m.delimiterTokens, token.Type) || m.whitespaceTokens[token.Start]
}

func (m *LexerModel) OnSuggestionChanged(suggestion input.Suggestion[any]) {
	token := m.currentToken(m.Value(), m.Cursor())

	if m.IsDelimiterToken(token) {
		if m.Cursor() < len(m.Value()) {
			token = m.tokens[token.Index-1]
			if m.IsDelimiterToken(token) {
				token = parser.Token{Start: token.End()}
			}
		} else {
			token = parser.Token{
				Start: token.End(),
			}
		}
	}
	m.setSelectedToken(&token)
	value := m.Value()
	newVal := value[:token.Start] + suggestion.Text
	if token.End() < len(value) {
		newVal += value[token.End():]
	}
	m.SetValue(newVal)
	m.SetCursor(token.Start + len(suggestion.Text) - suggestion.CursorOffset)

}

func (m *LexerModel) OnSuggestionUnselected() {
	m.setSelectedToken(nil)
}

func (m *LexerModel) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	return false
}

func (m *LexerModel) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	return true
}

func (m *LexerModel) CompletableTokenBeforeCursor() string {
	token := m.CurrentToken()
	if m.IsDelimiterToken(token) {
		// Don't filter suggestions on delimiters
		return ""
	}
	return m.currentTokenBeforeCursor(token)
}

func (m *LexerModel) CurrentTokenBeforeCursor() string {
	token := m.CurrentToken()
	return m.currentTokenBeforeCursor(token)
}

func (m *LexerModel) currentTokenBeforeCursor(token parser.Token) string {
	start := token.Start
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
