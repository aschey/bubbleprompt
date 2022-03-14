package prompt

import (
	"regexp"
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
	previousCommands        []string
	executorModel           *tea.Model
	modelState              modelState
	delimiterRegex          *regexp.Regexp
	lastTypedCursorPosition int
	typedText               string
	ready                   bool
	err                     error
}

func New(completer Completer, executor Executor, opts ...Option) Model {
	textInput := commandinput.New()
	textInput.Focus()

	model := Model{
		completer:      newCompleterModel(completer),
		executor:       executor,
		textInput:      textInput,
		delimiterRegex: regexp.MustCompile(`\s+`),
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
	return filterHasPrefix(search, suggestions,
		func(s Suggestion) string { return s.Text })
}

func FilterCompletionTextHasPrefix(search string, suggestions Suggestions) Suggestions {
	return filterHasPrefix(search, suggestions,
		func(s Suggestion) string { return s.CompletionText })
}

func filterHasPrefix(search string, suggestions Suggestions,
	textFunc func(s Suggestion) string) Suggestions {
	cleanedSearch := strings.TrimSpace(strings.ToLower(search))
	filtered := []Suggestion{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(textFunc(s)), cleanedSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

func (m *Model) SetPrompt(prompt string) {
	m.textInput.Prompt = prompt
}

func (m *Model) SetDelimiterRegex(delimiterRegex string) error {
	regex, err := regexp.Compile(delimiterRegex)
	if err != nil {
		return err
	}
	m.delimiterRegex = regex
	m.textInput.SetDelimiterRegex(delimiterRegex)
	return nil
}

func (m *Model) SetStringRegex(stringRegex string) error {
	m.textInput.SetStringRegex(stringRegex)
	return nil
}

func (m Model) CommandCompleted() bool {
	return m.textInput.CommandCompleted()
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
