package parserinput

import (
	"github.com/alecthomas/chroma/v2"
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
	styleLexer       chroma.Lexer
	style            *chroma.Style
	tokens           []parser.Token
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

func (m *LexerModel) createWhitespaceToken(start int, end int) parser.Token {
	token := parser.Token{
		Value: m.Value()[start:end],
		Start: start,
	}
	m.whitespaceTokens[start] = true
	return token
}

func (m *LexerModel) updateTokens() error {
	tokens, err := m.lexer.Lex(m.Value())
	if err != nil {
		return err
	}
	fullTokens := []parser.Token{}
	m.whitespaceTokens = make(map[int]bool)
	last := 0
	for i, token := range tokens {
		if i > 0 {
			prevEnd := tokens[i-1].End()
			if prevEnd < token.Start {
				// This part of the input was ignored by the lexer
				// so insert a dummy token to account for it
				fullTokens = append(fullTokens, m.createWhitespaceToken(prevEnd, token.Start))
			}
		}
		fullTokens = append(fullTokens, token)
		last = token.End()
	}

	// Check for trailing whitespace
	if m.Cursor() > last {
		fullTokens = append(fullTokens, m.createWhitespaceToken(last, m.Cursor()))
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
		// It will get reset during OnUpdateFinish
		if msg.Type != tea.KeyEnter {
			m.err = err
		}
	}

	return cmd
}

func (m *LexerModel) Error() error {
	return m.err
}

func (m *LexerModel) View(viewMode input.ViewMode) string {
	if m.styleLexer == nil {
		viewBuilder := input.NewViewBuilder(m.Cursor(), lipgloss.NewStyle(), " ", !m.textinput.Blink())
		viewBuilder.Render(m.Value(), 0, lipgloss.NewStyle())
		return m.prompt + viewBuilder.GetView()
	}
	iter, err := m.styleLexer.Tokenise(nil, m.textinput.Value())
	if err != nil {
		println(err)
	}
	return m.prompt + m.inputFormatter(iter, viewMode)
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
	return m.Cursor() == token.End() && tokenStr == suggestion.Text
}

func (m *LexerModel) currentToken(text string, tokenPos int) (int, parser.Token) {
	if len(m.tokens) == 0 {
		return -1, parser.Token{}
	}
	if len(m.Value()) == 0 {
		return 0, m.tokens[0]
	}

	for i, token := range m.tokens {
		if i == len(m.tokens)-1 || (tokenPos >= token.Start && tokenPos < token.End()) {
			return i, token
		}
	}
	return -1, parser.Token{}
}

func (m *LexerModel) CurrentToken() (int, parser.Token) {
	return m.currentToken(m.Value(), m.Cursor()-1)
}

func (m *LexerModel) FindLast(filter func(token parser.Token, symbol string) bool) *parser.Token {
	currentIndex, _ := m.CurrentToken()
	for i := currentIndex; i >= 0; i-- {
		token := m.tokens[i]

		if filter(token, token.Type) {
			return &token
		}
	}

	return nil
}

func (m *LexerModel) PreviousToken() (int, *parser.Token) {
	index, _ := m.CurrentToken()
	if index <= 0 {
		return -1, nil
	}
	return index - 1, &m.tokens[index-1]
}

func (m *LexerModel) CompletionText(text string) string {
	_, token := m.currentToken(text, m.Cursor()-1)
	return token.Value
}

func (m *LexerModel) Tokens() []parser.Token {
	return m.tokens
}

func (m *LexerModel) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[any]) tea.Cmd {
	return nil
}

func (m *LexerModel) IsDelimiterToken(token parser.Token) bool {
	// Dummy whitespace tokens won't be registered with the lexer so check them separately
	return slices.Contains(m.delimiters, token.Value) || slices.Contains(m.delimiterTokens, token.Type) || m.whitespaceTokens[token.Start]
}

func (m *LexerModel) OnSuggestionChanged(suggestion input.Suggestion[any]) {
	i, token := m.currentToken(m.Value(), m.Cursor())

	if m.IsDelimiterToken(token) {
		if m.Cursor() < len(m.Value()) {
			token = m.tokens[i-1]
			if m.IsDelimiterToken(token) {
				token = parser.Token{Start: token.End()}
			}
		} else {
			token = parser.Token{
				Start: token.End(),
			}
		}
	}

	value := m.Value()
	newVal := value[:token.Start] + suggestion.Text
	if token.End() < len(value) {
		newVal += value[token.End():]
	}
	m.SetValue(newVal)
	m.SetCursor(token.Start + len(suggestion.Text) - suggestion.CursorOffset)

}

func (m *LexerModel) OnSuggestionUnselected() {}

func (m *LexerModel) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	return false
}

func (m *LexerModel) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	return true
}

func (m *LexerModel) CompletableTokenBeforeCursor() string {
	_, token := m.CurrentToken()
	if m.IsDelimiterToken(token) {
		// Don't filter suggestions on delimiters
		return ""
	}
	return m.currentTokenBeforeCursor(token)
}

func (m *LexerModel) CurrentTokenBeforeCursor() string {
	_, token := m.CurrentToken()

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
