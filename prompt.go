package prompt

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type modelState int

const (
	completing modelState = iota
	executing
)

type Executor func(input string, selected *input.Suggestion, suggestions []input.Suggestion) tea.Model

type Model struct {
	completer               completerModel
	executor                Executor
	textInput               input.Input
	viewport                viewport.Model
	Formatters              input.Formatters
	previousCommands        []string
	executorModel           *tea.Model
	modelState              modelState
	scrollbar               string
	scrollbarThumb          string
	lastTypedCursorPosition int
	typedText               string
	ready                   bool
	err                     error
}

func New(completer Completer, executor Executor, textInput input.Input, opts ...Option) Model {
	model := Model{
		completer: newCompleterModel(completer, 6),
		executor:  executor,
		textInput: textInput,
		Formatters: input.Formatters{
			Name: input.SuggestionText{
				SelectedStyle: lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("240")).
					Background(lipgloss.Color("14")),
				Style: lipgloss.NewStyle().Background(lipgloss.Color("14")),
			},
			Description: input.SuggestionText{
				SelectedStyle: lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("240")).
					Background(lipgloss.Color("37")),
				Style: lipgloss.NewStyle().Background(lipgloss.Color("37")),
			},
			DefaultPlaceholder: input.Text{
				Style: lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
			},
			SelectedSuggestion: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		},
	}
	model.SetScrollbarColor(lipgloss.Color("14"))
	model.SetScrollbarThumbColor(lipgloss.Color("240"))
	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
}

func (m *Model) SetScrollbarColor(color lipgloss.TerminalColor) {
	m.scrollbar = lipgloss.NewStyle().Background(color).Render(" ")
}

func (m *Model) SetScrollbarThumbColor(color lipgloss.TerminalColor) {
	m.scrollbarThumb = lipgloss.NewStyle().Background(color).Render(" ")
}

func (m *Model) SetMaxSuggestions(maxSuggestions int) {
	m.completer.maxSuggestions = maxSuggestions
}

func FilterHasPrefix(search string, suggestions []input.Suggestion) []input.Suggestion {
	return filterHasPrefix(search, suggestions,
		func(s input.Suggestion) string { return s.Text })
}

func FilterCompletionTextHasPrefix(search string, suggestions []input.Suggestion) []input.Suggestion {
	return filterHasPrefix(search, suggestions,
		func(s input.Suggestion) string { return s.CompletionText })
}

func filterHasPrefix(search string, suggestions []input.Suggestion,
	textFunc func(s input.Suggestion) string) []input.Suggestion {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))
	filtered := []input.Suggestion{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(textFunc(s)), cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.completer.Init())
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}
