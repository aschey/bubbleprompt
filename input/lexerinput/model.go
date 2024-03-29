package lexerinput

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"golang.org/x/exp/slices"
)

type Model[T any] struct {
	textinput         textinput.Model
	lexer             parser.Lexer
	tokenFormatter    parser.Formatter
	formatters        Formatters
	selectedToken     *input.Token
	tokens            []input.Token
	formatterTokens   []parser.FormatterToken
	delimiterTokens   []string
	delimiters        []string
	whitespaceTokens  map[int]bool
	prompt            string
	currentSuggestion *string
	err               error
}

func NewModel[T any](lexer parser.Lexer, options ...Option[T]) *Model[T] {
	textinput := textinput.New()

	model := &Model[T]{
		lexer:            lexer,
		textinput:        textinput,
		prompt:           "> ",
		tokens:           []input.Token{},
		formatterTokens:  []parser.FormatterToken{},
		formatters:       DefaultFormatters(),
		whitespaceTokens: make(map[int]bool),
	}
	for _, option := range options {
		option(model)
	}
	return model
}

func (m *Model[T]) createWhitespaceToken(start int, end int, index int) input.Token {
	token := input.Token{
		Value: string(m.Runes()[start:end]),
		Start: start,
		Index: index,
	}
	m.whitespaceTokens[start] = true
	return token
}

func (m *Model[T]) updateTokens() error {
	tokens, err := m.lexer.Lex(m.Value())
	if err != nil {
		return err
	}

	if m.tokenFormatter != nil {
		m.formatterTokens, err = m.tokenFormatter.Lex(m.Value(), m.selectedToken)
		if err != nil {
			return err
		}
	}

	fullTokens := []input.Token{}
	m.whitespaceTokens = make(map[int]bool)
	last := 0
	index := 0
	for i, token := range tokens {
		if i > 0 {
			prevEnd := tokens[i-1].End()
			if prevEnd < token.Start {
				// This part of the input was ignored by the lexer
				// so insert a dummy token to account for it
				fullTokens = append(
					fullTokens,
					m.createWhitespaceToken(prevEnd, token.Start, index),
				)
				index++
			}
		}
		token.Index = index
		index++
		fullTokens = append(fullTokens, token)
		last = token.End()
	}

	// Check for trailing whitespace
	if m.CursorIndex() > last {
		fullTokens = append(
			fullTokens,
			m.createWhitespaceToken(last, len([]rune(m.textinput.Value())), index),
		)
	}
	m.tokens = fullTokens

	return nil
}

func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
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

func (m Model[T]) Error() error {
	return m.err
}

func (m Model[T]) unstyledView(text []rune, showCursor bool, viewMode input.ViewMode) string {
	viewBuilder := input.NewViewBuilder(m.CursorIndex(), m.formatters.Cursor, " ", showCursor)
	viewBuilder.Render(text, 0, lipgloss.NewStyle())
	return m.renderWithPlaceholder(viewBuilder, viewMode)
}

func (m Model[T]) styledView(
	formatterTokens []parser.FormatterToken,
	showCursor bool,
	viewMode input.ViewMode,
) string {
	viewBuilder := input.NewViewBuilder(m.CursorIndex(), m.formatters.Cursor, " ", showCursor)
	for _, token := range formatterTokens {
		viewBuilder.Render(
			[]rune(strings.TrimRight(token.Value, "\n")),
			viewBuilder.ViewLen(),
			token.Style,
		)
	}
	return m.renderWithPlaceholder(viewBuilder, viewMode)
}

func (m Model[T]) renderWithPlaceholder(
	viewBuilder *input.ViewBuilder,
	viewMode input.ViewMode,
) string {
	if m.currentSuggestion != nil && viewMode == input.Interactive {
		current := m.CurrentToken()
		suggestionRunes := []rune(*m.currentSuggestion)
		if current.Index > -1 {
			value := current.Value
			if m.IsDelimiterToken(current) {
				if current.Index < len(m.Tokens())-1 {
					// Cursor is before the last token, don't render placeholder
					return m.prompt + viewBuilder.View()
				}
				// Current token is a delimiter, don't try to filter it on the prefix
				value = ""
			}
			suggestionRunes = suggestionRunes[len([]rune(value)):]
			// Render placeholder only if the prefix matches
			if strings.HasPrefix(*m.currentSuggestion, value) {
				viewBuilder.RenderPlaceholder(
					suggestionRunes,
					viewBuilder.ViewLen(),
					m.formatters.Placeholder,
				)
			}

		} else {
			// Nothing typed yet, always render the placeholder
			viewBuilder.RenderPlaceholder(suggestionRunes, viewBuilder.ViewLen(), m.formatters.Placeholder)
		}

	}
	return m.prompt + viewBuilder.View()
}

func (m Model[T]) View(viewMode input.ViewMode) string {
	showCursor := !m.textinput.Cursor.Blink
	if viewMode == input.Static {
		showCursor = false
	}
	if m.tokenFormatter == nil {
		return m.unstyledView(m.Runes(), showCursor, viewMode)
	}

	return m.styledView(m.formatterTokens, showCursor, viewMode)
}

func (m Model[T]) FormatText(text string) string {
	if m.tokenFormatter == nil {
		return m.unstyledView([]rune(text), false, input.Static)
	}
	formatterTokens, _ := m.tokenFormatter.Lex(text, nil)
	return m.styledView(formatterTokens, false, input.Static)
}

func (m *Model[T]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m Model[T]) Focused() bool {
	return m.textinput.Focused()
}

func (m Model[T]) Value() string {
	return m.textinput.Value()
}

func (m Model[T]) Runes() []rune {
	return []rune(m.textinput.Value())
}

func (m *Model[T]) ResetValue() {
	m.textinput.SetValue("")
	_ = m.updateTokens()
}

func (m *Model[T]) SetValue(value string) {
	m.textinput.SetValue(value)
	m.err = m.updateTokens()
}

func (m *Model[T]) setSelectedToken(token *input.Token) {
	m.selectedToken = token
	_ = m.updateTokens()
}

func (m *Model[T]) Blur() {
	m.textinput.Blur()
}

func (m Model[T]) CursorIndex() int {
	return m.textinput.Position()
}

func (m Model[T]) CursorOffset() int {
	cursorIndex := m.CursorIndex()
	runesBeforeCursor := m.Runes()[:cursorIndex]
	return runewidth.StringWidth(string(runesBeforeCursor))
}

func (m *Model[T]) SetCursor(cursor int) {
	m.textinput.SetCursor(cursor)
}

func (m *Model[T]) SetCursorMode(cursorMode cursor.Mode) tea.Cmd {
	return m.textinput.Cursor.SetMode(cursorMode)
}

// Formatters returns the formatters used by the input.
func (m Model[T]) Formatters() Formatters {
	return m.formatters
}

func (m *Model[T]) SetFormatters(formatters Formatters) {
	m.formatters = formatters
}

func (m Model[T]) Prompt() string {
	return m.prompt
}

func (m *Model[T]) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m Model[T]) ShouldSelectSuggestion(suggestion suggestion.Suggestion[T]) bool {
	token := m.CurrentToken()
	tokenStr := token.Value
	return m.CursorIndex() == token.End() && tokenStr == suggestion.Text
}

func (m Model[T]) currentToken(
	runes []rune,
	tokenPos int,
	roundingBehavior input.RoundingBehavior,
) input.Token {
	return input.FindCurrentToken(
		runes,
		m.tokens,
		m.CursorIndex(),
		roundingBehavior,
		func(s string, last input.Token) bool {
			return m.IsDelimiterToken(last)
		},
	)
}

func (m Model[T]) CurrentToken() input.Token {
	return m.currentToken(m.Runes(), m.CursorIndex()-1, input.RoundUp)
}

func (m Model[T]) CurrentTokenRoundDown() input.Token {
	return m.currentToken(m.Runes(), m.CursorIndex()-1, input.RoundDown)
}

func (m Model[T]) FindLast(filter func(token input.Token, symbol string) bool) *input.Token {
	currentToken := m.CurrentToken()
	for i := currentToken.Index; i >= 0; i-- {
		token := m.tokens[i]

		if filter(token, token.Type) {
			return &token
		}
	}

	return nil
}

func (m Model[T]) PreviousToken() *input.Token {
	currentToken := m.CurrentToken()
	if currentToken.Index <= 0 {
		return nil
	}
	return &m.tokens[currentToken.Index-1]
}

func (m Model[T]) SuggestionRunes(runes []rune) []rune {
	token := m.currentToken(runes, m.CursorIndex()-1, input.RoundUp)
	return []rune(token.Value)
}

func (m Model[T]) Tokens() []input.Token {
	return m.tokens
}

func (m Model[T]) TokensBeforeCursor() []input.Token {
	tokens := []input.Token{}
	cursor := m.CursorIndex()
	for _, token := range m.tokens {
		if token.End() <= cursor {
			tokens = append(tokens, token)
		} else {
			tokens = append(tokens, input.Token{
				Value: string([]rune(token.Value)[:cursor-token.Start]),
				Start: token.Start,
				Index: token.Index,
				Type:  token.Type,
			})
			break
		}
	}
	return tokens
}

func (m Model[T]) TokenValues() []string {
	tokens := []string{}
	for _, token := range m.tokens {
		tokens = append(tokens, token.Value)
	}
	return tokens
}

func (m *Model[T]) OnUpdateFinish(
	msg tea.Msg,
	suggestion *suggestion.Suggestion[T],
	isSelected bool,
) tea.Cmd {
	if suggestion == nil || isSelected {
		m.currentSuggestion = nil
	} else if !isSelected {
		m.currentSuggestion = &suggestion.Text
	}
	return nil
}

func (m Model[T]) IsDelimiterToken(token input.Token) bool {
	// Dummy whitespace tokens won't be registered with the lexer so check them separately
	return slices.Contains(m.delimiters, token.Value) ||
		slices.Contains(m.delimiterTokens, token.Type) ||
		m.whitespaceTokens[token.Start]
}

func (m *Model[T]) OnSuggestionChanged(suggestion suggestion.Suggestion[T]) {
	runes := m.Runes()
	token := m.currentToken(runes, m.CursorIndex(), input.RoundUp)

	if m.IsDelimiterToken(token) {
		if m.CursorIndex() < len(runes) && token.Index < len(m.tokens)-1 &&
			!m.IsDelimiterToken(m.tokens[token.Index+1]) {
			token = m.tokens[token.Index+1]
		} else {
			token = input.Token{
				Start: token.End(),
			}
		}
	}
	m.setSelectedToken(&token)

	suggestionRunes := []rune(suggestion.Text)
	newVal := append(m.Runes()[:token.Start], suggestionRunes...)
	if token.End() < len(runes) {
		newVal = append(newVal, runes[token.End():]...)
	}
	m.SetValue(string(newVal))
	m.SetCursor(token.Start + len(suggestionRunes) - suggestion.CursorOffset)
}

func (m *Model[T]) OnSuggestionUnselected() {
	m.setSelectedToken(nil)
}

func (m Model[T]) ShouldClearSuggestions(prevText []rune, msg tea.KeyMsg) bool {
	return false
}

func (m Model[T]) ShouldUnselectSuggestion(prevText []rune, msg tea.KeyMsg) bool {
	return true
}

func (m Model[T]) CompletableTokenBeforeCursor() string {
	token := m.CurrentToken()
	if m.IsDelimiterToken(token) {
		// Don't filter suggestions on delimiters
		return ""
	}
	return string(m.currentTokenBeforeCursor(token))
}

func (m Model[T]) CurrentTokenBeforeCursor() string {
	token := m.CurrentToken()
	return string(m.currentTokenBeforeCursor(token))
}

func (m Model[T]) currentTokenBeforeCursor(token input.Token) []rune {
	start := token.Start
	cursor := m.CursorIndex()
	if start > cursor {
		return []rune("")
	}
	val := m.Runes()[start:cursor]
	return val
}

func (m *Model[T]) OnExecutorFinished() {
	// Clear out error once input text is reset
	m.err = nil
}
