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

type Model[T any] struct {
	lexerModel *parserinput.LexerModel[T]
}

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

func (m *Model[T]) CurrentToken() editor.Token {
	return m.lexerModel.CurrentToken()
}

func (m *Model[T]) CurrentTokenBeforeCursor() string {
	return m.lexerModel.CompletableTokenBeforeCursor()
}

func (m *Model[T]) TokenValues() []string {
	tokenValues := []string{}
	tokens := m.Tokens()
	for _, token := range tokens {
		tokenValues = append(tokenValues, token.Value)
	}
	return tokenValues
}

func (m *Model[T]) AllTokens() []editor.Token {
	return m.lexerModel.Tokens()
}

func (m *Model[T]) Tokens() []editor.Token {
	return m.filterWhitespaceTokens(m.lexerModel.Tokens())
}

func (m *Model[T]) AllTokensBeforeCursor() []editor.Token {
	return m.lexerModel.Tokens()
}

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

func (m *Model[T]) Init() tea.Cmd {
	return m.lexerModel.Init()
}

func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	return m.lexerModel.OnUpdateStart(msg)
}

func (m *Model[T]) View(viewMode editor.ViewMode) string {
	return m.lexerModel.View(viewMode)
}

func (m *Model[T]) Focus() tea.Cmd {
	return m.lexerModel.Focus()
}

func (m *Model[T]) Focused() bool {
	return m.lexerModel.Focused()
}

func (m *Model[T]) Value() string {
	return m.lexerModel.Value()
}

func (m *Model[T]) Runes() []rune {
	return m.lexerModel.Runes()
}

func (m *Model[T]) ResetValue() {
	m.lexerModel.ResetValue()
}

func (m *Model[T]) SetValue(value string) {
	m.lexerModel.SetValue(value)
}

func (m *Model[T]) Blur() {
	m.lexerModel.Blur()
}

func (m *Model[T]) CursorOffset() int {
	return m.lexerModel.CursorOffset()
}

func (m *Model[T]) CursorIndex() int {
	return m.lexerModel.CursorIndex()
}

func (m *Model[T]) SetCursor(cursor int) {
	m.lexerModel.SetCursor(cursor)
}

func (m *Model[T]) SetCursorMode(cursorMode textinput.CursorMode) tea.Cmd {
	return m.lexerModel.SetCursorMode(cursorMode)
}

func (m *Model[T]) Prompt() string {
	return m.lexerModel.Prompt()
}

func (m *Model[T]) SetPrompt(prompt string) {
	m.lexerModel.SetPrompt(prompt)
}

func (m *Model[T]) ShouldSelectSuggestion(suggestion editor.Suggestion[T]) bool {
	return m.lexerModel.ShouldSelectSuggestion(suggestion)
}

func (m *Model[T]) CompletionRunes(runes []rune) []rune {
	return m.lexerModel.CompletionRunes(runes)
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
