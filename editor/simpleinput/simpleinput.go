package simpleinput

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/editor/parser"
	"github.com/aschey/bubbleprompt/editor/parserinput"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A Model is a simple editor for handling simple token-based inputs without any special parsing required.
type Model[T any] struct {
	lexerModel *parserinput.LexerModel[T]
}

// New creates new a model.
func New[T any](options ...Option[T]) *Model[T] {
	settings := &settings[T]{
		delimiterRegex:    `\s+`,
		tokenRegex:        `("[^"]*"?)|('[^']*'?)|[^\s]+`,
		selectedTextStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		lexerOptions:      []parserinput.Option[T]{},
	}
	for _, option := range options {
		if err := option(settings); err != nil {
			panic(err)
		}
	}
	lexerDefinition := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Token", Pattern: settings.tokenRegex},
		{Name: "Delimiter", Pattern: settings.delimiterRegex},
	})

	var formatter parser.Formatter
	if settings.formatter != nil {
		formatter = *settings.formatter
	} else {
		formatter = simpleFormatter{
			lexer:             lexerDefinition,
			selectedTextStyle: settings.selectedTextStyle,
		}
	}

	lexer := parser.NewParticipleLexer(lexerDefinition)

	m := &Model[T]{
		parserinput.NewLexerModel(lexer,
			append(settings.lexerOptions,
				parserinput.WithDelimiterTokens[T]("Delimiter"),
				parserinput.WithFormatter[T](formatter),
			)...),
	}

	return m
}

// CurrentToken returns the token under the cursor.
func (m *Model[T]) CurrentToken() editor.Token {
	return m.lexerModel.CurrentToken()
}

// CurrentTokenBeforeCursor returns the portion of the token under the cursor
// that comes before the cursor position.
func (m *Model[T]) CurrentTokenBeforeCursor() string {
	return m.lexerModel.CompletableTokenBeforeCursor()
}

// TokenValues returns the tokenized input text.
// This **does not** include delimiter tokens.
func (m *Model[T]) TokenValues() []string {
	tokenValues := []string{}
	tokens := m.Tokens()
	for _, token := range tokens {
		tokenValues = append(tokenValues, token.Value)
	}
	return tokenValues
}

// AllTokens returns the tokenized input.
// This **does** include delimiter tokens.
func (m *Model[T]) AllTokens() []editor.Token {
	return m.lexerModel.Tokens()
}

// Tokens returns the tokenized input.
// This **does not** include delimiter tokens.
func (m *Model[T]) Tokens() []editor.Token {
	return m.filterWhitespaceTokens(m.lexerModel.Tokens())
}

// AllTokensBeforeCursor returns the tokenized input up to the cursor position.
// This **does not** include delimiter tokens.
func (m *Model[T]) AllTokensBeforeCursor() []editor.Token {
	return m.lexerModel.Tokens()
}

// AllTokensBeforeCursor returns the tokenized input up to the cursor position.
// This **does** include delimiter tokens.
func (m *Model[T]) TokensBeforeCursor() []editor.Token {
	return m.filterWhitespaceTokens(m.lexerModel.TokensBeforeCursor())
}

func (m *Model[T]) filterWhitespaceTokens(allTokens []editor.Token) []editor.Token {
	tokens := []editor.Token{}
	for _, token := range allTokens {
		if token.Type == "Token" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

// Init is part of the editor interface.
// It does not need to be invoked by end users.
func (m *Model[T]) Init() tea.Cmd {
	return m.lexerModel.Init()
}

// OnUpdateStart is part of the editor interface.
// It does not need to be invoked by end users.
func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	return m.lexerModel.OnUpdateStart(msg)
}

// View is part of the editor interface.
// It does not need to be invoked by end users.
func (m *Model[T]) View(viewMode editor.ViewMode) string {
	return m.lexerModel.View(viewMode)
}

// Focus sets the keyboard focus on the editor so the user can enter text.
func (m *Model[T]) Focus() tea.Cmd {
	return m.lexerModel.Focus()
}

// Focused returns whether the keyboard is focused on the editor.
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

// ResetValue clears the editor.
func (m *Model[T]) ResetValue() {
	m.lexerModel.ResetValue()
}

// SetValue sets the text of the editor.
func (m *Model[T]) SetValue(value string) {
	m.lexerModel.SetValue(value)
}

// Blur removes the focus from the editor.
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
func (m *Model[T]) SetCursorMode(cursorMode textinput.CursorMode) tea.Cmd {
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

func (m *Model[T]) ShouldSelectSuggestion(suggestion editor.Suggestion[T]) bool {
	return m.lexerModel.ShouldSelectSuggestion(suggestion)
}

func (m *Model[T]) SuggestionRunes(runes []rune) []rune {
	return m.lexerModel.SuggestionRunes(runes)
}

func (m *Model[T]) OnUpdateFinish(msg tea.Msg, suggestion *editor.Suggestion[T], isSelected bool) tea.Cmd {
	return m.lexerModel.OnUpdateFinish(msg, suggestion, isSelected)
}

func (m *Model[T]) OnSuggestionChanged(suggestion editor.Suggestion[T]) {
	m.lexerModel.OnSuggestionChanged(suggestion)
}

func (m *Model[T]) OnExecutorFinished() {
	m.lexerModel.OnExecutorFinished()
}

func (m *Model[T]) OnSuggestionUnselected() {
	m.lexerModel.OnSuggestionUnselected()
}

func (m *Model[T]) ShouldClearSuggestions(prevText []rune, msg tea.KeyMsg) bool {
	return m.lexerModel.ShouldClearSuggestions(prevText, msg)
}

func (m *Model[T]) ShouldUnselectSuggestion(prevText []rune, msg tea.KeyMsg) bool {
	return m.lexerModel.ShouldUnselectSuggestion(prevText, msg)
}
