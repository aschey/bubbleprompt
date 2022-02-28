package prompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Formatters struct {
	Name               SuggestionText
	Description        SuggestionText
	DefaultPlaceholder Text
	SelectedSuggestion lipgloss.Style
}

type modelState int

const (
	completing modelState = iota
	executing
)

type Executor func(input string, selected *Suggestion, suggestions Suggestions) tea.Model

type Model struct {
	completer               completerModel
	executor                Executor
	textInput               textinput.Model
	viewport                viewport.Model
	Formatters              Formatters
	previousCommands        []string
	executorModel           *tea.Model
	modelState              modelState
	lastTypedCursorPosition int
	typedText               string
	listPosition            int
	ready                   bool
	err                     error
}

func New(completer Completer, executor Executor, opts ...Option) Model {
	textInput := textinput.New()
	textInput.Placeholder = "first-option"
	textInput.Focus()

	model := Model{
		completer: newCompleterModel(completer),
		executor:  executor,
		textInput: textInput,
		Formatters: Formatters{
			Name: SuggestionText{
				SelectedStyle: lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("240")).
					Background(lipgloss.Color("14")),
				Style: lipgloss.NewStyle().Background(lipgloss.Color("14")),
			},
			Description: SuggestionText{
				SelectedStyle: lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("240")).
					Background(lipgloss.Color("37")),
				Style: lipgloss.NewStyle().Background(lipgloss.Color("37")),
			},
			DefaultPlaceholder: Text{
				Style: lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
			},
			SelectedSuggestion: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		},
		listPosition: -1,
	}

	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
}

func FilterHasPrefix(search string, suggestions Suggestions) Suggestions {
	lowerSearch := strings.ToLower(search)
	filtered := []Suggestion{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(s.Name), lowerSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

func (m *Model) SetPrompt(prompt string) {
	m.textInput.Prompt = prompt
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

func (m *Model) unselectSuggestion() {
	m.listPosition = -1
}

func (m Model) isSuggestionSelected() bool {
	return m.listPosition > -1
}

func (m Model) getSelectedSuggestion() *Suggestion {
	if m.isSuggestionSelected() {
		return &m.completer.suggestions[m.listPosition]
	}
	return nil
}
