package simpleinput

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parser"
	"github.com/aschey/bubbleprompt/input/parserinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	lexerModel *parserinput.LexerModel
}

func New(options ...Option) *Model {
	settings := &settings{
		delimiterRegex:    `\s+`,
		tokenRegex:        `[^\s]+`,
		selectedTextStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
	}
	for _, option := range options {
		if err := option(settings); err != nil {
			panic(err)
		}
	}
	lexerDefinition := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Delimiter", Pattern: settings.delimiterRegex},
		{Name: "Token", Pattern: settings.tokenRegex},
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

	m := &Model{
		parserinput.NewLexerModel(lexer, parserinput.WithDelimiterTokens("Delimiter"), parserinput.WithFormatter(formatter)),
	}

	return m
}

func (m *Model) CurrentToken() parser.Token {
	return m.lexerModel.CurrentToken()
}

func (m *Model) CurrentTokenBeforeCursor() string {
	return m.lexerModel.CompletableTokenBeforeCursor()
}

func (m *Model) TokenValues() []string {
	tokenValues := []string{}
	tokens := m.Tokens()
	for _, token := range tokens {
		tokenValues = append(tokenValues, token.Value)
	}
	return tokenValues
}

func (m *Model) AllTokens() []parser.Token {
	return m.lexerModel.Tokens()
}

func (m *Model) Tokens() []parser.Token {
	tokens := []parser.Token{}
	allTokens := m.lexerModel.Tokens()
	for _, token := range allTokens {
		if token.Type == "Token" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func (m *Model) Init() tea.Cmd {
	return m.lexerModel.Init()
}

func (m *Model) OnUpdateStart(msg tea.Msg) tea.Cmd {
	return m.lexerModel.OnUpdateStart(msg)
}

func (m *Model) View(viewMode input.ViewMode) string {
	return m.lexerModel.View(viewMode)
}

func (m *Model) Focus() tea.Cmd {
	return m.lexerModel.Focus()
}

func (m *Model) Focused() bool {
	return m.lexerModel.Focused()
}

func (m *Model) Value() string {
	return m.lexerModel.Value()
}

func (m *Model) ResetValue() {
	m.lexerModel.ResetValue()
}

func (m *Model) SetValue(value string) {
	m.lexerModel.SetValue(value)
}

func (m *Model) Blur() {
	m.lexerModel.Blur()
}

func (m *Model) Cursor() int {
	return m.lexerModel.Cursor()
}

func (m *Model) SetCursor(cursor int) {
	m.lexerModel.SetCursor(cursor)
}

func (m *Model) Prompt() string {
	return m.lexerModel.Prompt()
}

func (m *Model) SetPrompt(prompt string) {
	m.lexerModel.SetPrompt(prompt)
}

func (m *Model) ShouldSelectSuggestion(suggestion input.Suggestion[any]) bool {
	return m.lexerModel.ShouldSelectSuggestion(suggestion)
}

func (m *Model) CompletionText(text string) string {
	return m.lexerModel.CompletionText(text)
}

func (m *Model) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[any], isSelected bool) tea.Cmd {
	return m.lexerModel.OnUpdateFinish(msg, suggestion, isSelected)
}

func (m *Model) OnSuggestionChanged(suggestion input.Suggestion[any]) {
	m.lexerModel.OnSuggestionChanged(suggestion)
}

func (m *Model) OnExecutorFinished() {
	m.lexerModel.OnExecutorFinished()
}

func (m *Model) OnSuggestionUnselected() {
	m.lexerModel.OnSuggestionUnselected()
}

func (m *Model) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	return m.lexerModel.ShouldClearSuggestions(prevText, msg)
}

func (m *Model) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	return m.lexerModel.ShouldUnselectSuggestion(prevText, msg)
}
