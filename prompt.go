package prompt

import (
	"strings"

	"github.com/aschey/bubbleprompt/commandinput"
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
	textInput               commandinput.Model
	viewport                viewport.Model
	Formatters              Formatters
	Separators              []string
	previousCommands        []string
	executorModel           *tea.Model
	modelState              modelState
	lastTypedCursorPosition int
	typedText               string
	ready                   bool
	err                     error
}

func New(completer Completer, executor Executor, opts ...Option) Model {
	textInput := commandinput.New()
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
	}

	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
}

func FilterHasPrefix(search string, suggestions Suggestions) Suggestions {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))
	filtered := []Suggestion{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(s.Name), cleanedSearch) {
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
