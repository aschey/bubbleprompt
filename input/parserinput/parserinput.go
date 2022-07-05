package parserinput

import (
	"github.com/alecthomas/participle/v2"
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Grammar interface {
	CurrentToken() string
}

type Model[T Grammar] struct {
	textinput  textinput.Model
	parser     *participle.Parser[T]
	parsedText *T
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

func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	expr, err := m.parser.ParseString("", m.Value())
	if err == nil {
		m.parsedText = expr
	}
	return cmd
}

func (m *Model[T]) View() string {
	return m.textinput.View()
}

func (m *Model[T]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *Model[T]) Focused() bool {
	return m.textinput.Focused()
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
	return true
}

func (m *Model[T]) CompletionText(text string) string {
	expr, _ := m.parser.ParseString("", text)
	return (*expr).CurrentToken()
}

func (m *Model[T]) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[T]) tea.Cmd {
	return nil
}

func (m *Model[T]) OnSuggestionChanged(suggestion input.Suggestion[T]) {}

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
	return (*m.parsedText).CurrentToken()
}
