package prompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Suggestion struct {
	Name        string
	Description string
	Placeholder string
}

type errMsg error

type Completer func(input string) []Suggestion

type Executor func(input string, selected *Suggestion, suggestions []Suggestion) string

type Model struct {
	completer          Completer
	executor           Executor
	suggestions        []Suggestion
	textInput          textinput.Model
	viewport           viewport.Model
	previousCommands   []string
	Name               SuggestionText
	Description        SuggestionText
	Placeholder        Text
	SelectedSuggestion Text
	typedText          string
	prevText           string
	updating           bool
	listPosition       int
	placeholderValue   string
	ready              bool
	err                error
}

func New(completer Completer, executor Executor, opts ...Option) Model {
	textInput := textinput.New()
	textInput.Focus()

	model := Model{
		completer: completer,
		executor:  executor,
		textInput: textInput,
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
		suggestions:  completer(""),
		listPosition: -1,
	}

	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
}

func FilterHasPrefix(search string, suggestions []Suggestion) []Suggestion {
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
	return textinput.Blink
}

type completionMsg []Suggestion

func (m Model) updateCompletions() tea.Cmd {
	return func() tea.Msg {
		filtered := m.completer(m.textInput.Value())

		return completionMsg(filtered)
	}
}

func (m Model) render() string {
	maxNameLen := 0
	maxDescLen := 0

	// Determine longest name and description to calculate padding
	for _, s := range m.suggestions {
		if len(s.Name) > maxNameLen {
			maxNameLen = len(s.Name)
		}
		if len(s.Description) > maxDescLen {
			maxDescLen = len(s.Description)
		}
	}

	// Calculate left offset for suggestions
	padding := lipgloss.
		NewStyle().
		PaddingLeft(len(m.textInput.Prompt + m.typedText)).
		Render("")

	textView := m.textInput.View() + m.Placeholder.format(m.placeholderValue)

	// If an item is selected, parse out the text portion and apply formatting
	if m.listPosition > -1 {
		prompt := m.textInput.Prompt
		value := m.textInput.Value()
		formattedSuggestion := m.SelectedSuggestion.format(value)
		remainder := textView[len(prompt)+len(value):]
		textView = prompt + formattedSuggestion + remainder

	}

	prompts := append(m.previousCommands, textView)

	for i, s := range m.suggestions {
		selected := i == m.listPosition
		name := m.Name.format(s.Name, maxNameLen, selected)
		description := m.Description.format(s.Description, maxDescLen, selected)

		line := lipgloss.JoinHorizontal(lipgloss.Bottom, padding, name, description)
		prompts = append(prompts, line)
	}

	// Reserve height for prompts that were filtered out
	extraHeight := 5 - len(m.suggestions) - 1
	if extraHeight > 0 {
		extraLines := strings.Repeat("\n", extraHeight)
		prompts = append(prompts, extraLines)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, prompts...)
	return ret
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}
