// Package simpleinput provides an implementation of the [input.Input] interface.
// It should be used for basic cases without the need for structured or cli-style input
package simpleinput

import (
	participlelexer "github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/lexerinput"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A Model is an input for handling simple token-based inputs without any special parsing required.
type Model[T any] struct {
	lexerModel *lexerinput.Model[T]
}

// New creates new a model.
func New[T any](options ...Option[T]) *Model[T] {
	settings := &settings[T]{
		delimiterRegex:    `\s+`,
		tokenRegex:        `("[^"]*"?)|('[^']*'?)|[^\s]+`,
		selectedTextStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		lexerOptions:      []lexerinput.Option[T]{},
	}
	for _, option := range options {
		option(settings)
	}
	lexerDefinition := participlelexer.MustSimple([]participlelexer.SimpleRule{
		{Name: "Token", Pattern: settings.tokenRegex},
		{Name: "Delimiter", Pattern: settings.delimiterRegex},
	})
	participleLexer := parser.NewParticipleLexer(lexerDefinition)

	var formatter parser.Formatter
	if settings.formatterFunc != nil {
		formatter = settings.formatterFunc(participleLexer)
	} else {
		formatter = simpleinputFormatter{
			lexer:             participleLexer,
			selectedTextStyle: settings.selectedTextStyle,
		}
	}

	m := &Model[T]{
		lexerinput.NewModel(participleLexer,
			append(settings.lexerOptions,
				lexerinput.WithDelimiterTokens[T]("Delimiter"),
				lexerinput.WithFormatter[T](formatter),
			)...),
	}

	return m
}

// CurrentToken returns the token under the cursor.
func (m *Model[T]) CurrentToken() input.Token {
	return m.lexerModel.CurrentToken()
}

// CurrentTokenBeforeCursor returns the portion of the token under the cursor
// that comes before the cursor position.
func (m *Model[T]) CurrentTokenBeforeCursor() string {
	return m.lexerModel.CompletableTokenBeforeCursor()
}

// WordTokenValues returns the tokenized input text.
// This does not include delimiter tokens.
func (m *Model[T]) WordTokenValues() []string {
	tokenValues := []string{}
	tokens := m.WordTokens()
	for _, token := range tokens {
		tokenValues = append(tokenValues, token.Unquote())
	}
	return tokenValues
}

// Tokens returns the tokenized input.
// This includes delimiter tokens.
func (m *Model[T]) Tokens() []input.Token {
	return m.lexerModel.Tokens()
}

// WordTokens returns the tokenized input.
// This does not include delimiter tokens.
func (m *Model[T]) WordTokens() []input.Token {
	return m.filterWhitespaceTokens(m.lexerModel.Tokens())
}

// TokensBeforeCursor returns the tokenized input up to the cursor position.
// This includes delimiter tokens.
func (m *Model[T]) TokensBeforeCursor() []input.Token {
	return m.lexerModel.Tokens()
}

// WordTokensBeforeCursor returns the tokenized input up to the cursor position.
// This does not include delimiter tokens.
func (m *Model[T]) WordTokensBeforeCursor() []input.Token {
	return m.filterWhitespaceTokens(m.lexerModel.TokensBeforeCursor())
}

func (m *Model[T]) filterWhitespaceTokens(allTokens []input.Token) []input.Token {
	tokens := []input.Token{}
	for _, token := range allTokens {
		if token.Type == "Token" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

// Focus sets the keyboard focus on the editor so the user can enter text.
func (m *Model[T]) Focus() tea.Cmd {
	return m.lexerModel.Focus()
}

// Focused returns whether the keyboard is focused on the input.
func (m *Model[T]) Focused() bool {
	return m.lexerModel.Focused()
}

// Value returns the raw text entered by the user.
func (m *Model[T]) Value() string {
	return m.lexerModel.Value()
}

// Runes returns the raw text entered by the user as a list of runes.
// This is useful for indexing and length checks because doing these
// operations on strings does not work well with some unicode characters.
func (m *Model[T]) Runes() []rune {
	return m.lexerModel.Runes()
}

// ResetValue clears the input.
func (m *Model[T]) ResetValue() {
	m.lexerModel.ResetValue()
}

// SetValue sets the text of the input.
func (m *Model[T]) SetValue(value string) {
	m.lexerModel.SetValue(value)
}

// Blur removes the focus from the input.
func (m *Model[T]) Blur() {
	m.lexerModel.Blur()
}

// CursorOffset returns the visual offset of the cursor in terms
// of number of terminal cells. Use this for calculating visual dimensions
// such as input width/height.
func (m *Model[T]) CursorOffset() int {
	return m.lexerModel.CursorOffset()
}

// CursorIndex returns the cursor index in terms of number of unicode characters.
// Use this to calculate input lengths in terms of number of characters entered.
func (m *Model[T]) CursorIndex() int {
	return m.lexerModel.CursorIndex()
}

// Set cursor sets the cursor position.
func (m *Model[T]) SetCursor(cursor int) {
	m.lexerModel.SetCursor(cursor)
}

// SetCursorMode sets the mode of the cursor.
func (m *Model[T]) SetCursorMode(cursorMode cursor.Mode) tea.Cmd {
	return m.lexerModel.SetCursorMode(cursorMode)
}

// Prompt returns the terminal prompt.
func (m *Model[T]) Prompt() string {
	return m.lexerModel.Prompt()
}

// SetPrompt sets the terminal prompt.
func (m *Model[T]) SetPrompt(prompt string) {
	m.lexerModel.SetPrompt(prompt)
}

// Init is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) Init() tea.Cmd {
	return m.lexerModel.Init()
}

// OnUpdateStart is part of the Input interface.
// It  should not be invoked by end users.
func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	return m.lexerModel.OnUpdateStart(msg)
}

// View is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) View(viewMode input.ViewMode) string {
	return m.lexerModel.View(viewMode)
}

// ShouldSelectSuggestion is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) ShouldSelectSuggestion(suggestion suggestion.Suggestion[T]) bool {
	return m.lexerModel.ShouldSelectSuggestion(suggestion)
}

// SuggestionRunes is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) SuggestionRunes(runes []rune) []rune {
	return m.lexerModel.SuggestionRunes(runes)
}

// OnUpdateFinish is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) OnUpdateFinish(msg tea.Msg, suggestion *suggestion.Suggestion[T], isSelected bool) tea.Cmd {
	return m.lexerModel.OnUpdateFinish(msg, suggestion, isSelected)
}

// OnSuggestionChanged is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) OnSuggestionChanged(suggestion suggestion.Suggestion[T]) {
	m.lexerModel.OnSuggestionChanged(suggestion)
}

// OnExecutorFinished is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) OnExecutorFinished() {
	m.lexerModel.OnExecutorFinished()
}

// OnSuggestionUnselected is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) OnSuggestionUnselected() {
	m.lexerModel.OnSuggestionUnselected()
}

// ShouldClearSuggestions is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) ShouldClearSuggestions(prevText []rune, msg tea.KeyMsg) bool {
	return m.lexerModel.ShouldClearSuggestions(prevText, msg)
}

// ShouldUnselectSuggestion is part of the Input interface.
// It should not be invoked by users of this library.
func (m *Model[T]) ShouldUnselectSuggestion(prevText []rune, msg tea.KeyMsg) bool {
	return m.lexerModel.ShouldUnselectSuggestion(prevText, msg)
}
