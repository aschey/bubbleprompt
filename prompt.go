package prompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type Executor func(input string, selected *Suggestion, suggestions Suggestions) tea.Model

type Formatters struct {
	Name               SuggestionText
	Description        SuggestionText
	Placeholder        Text
	SelectedSuggestion Text
}

type Model struct {
	completer        completerModel
	executor         Executor
	textInput        textinput.Model
	viewport         viewport.Model
	Formatters       Formatters
	previousCommands []string
	executorModel    *tea.Model
	typedText        string
	listPosition     int
	placeholderValue string
	ready            bool
	err              error
}

func New(completer Completer, executor Executor, opts ...Option) Model {
	textInput := textinput.New()
	textInput.Focus()

	model := Model{
		completer: newCompleterModel(completer),
		executor:  executor,
		textInput: textInput,
		Formatters: Formatters{
			Name: SuggestionText{
				SelectedForegroundColor: "240",
				BackgroundColor:         "14",
				SelectedBackgroundColor: "14",
			},
			Description: SuggestionText{
				SelectedForegroundColor: "240",
				BackgroundColor:         "37",
				SelectedBackgroundColor: "37",
			},
			Placeholder: Text{
				ForegroundColor: "6",
			},
			SelectedSuggestion: Text{
				ForegroundColor: "10",
			},
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
