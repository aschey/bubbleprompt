package prompt

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Suggest struct {
	Name        string
	Description string
}

type errMsg error

type Model struct {
	suggestions         []Suggest
	filteredSuggestions []Suggest
	textInput           textinput.Model
	prevText            string
	updating            bool
	err                 error
}

func New(opts ...Option) Model {
	textInput := textinput.New()
	textInput.Focus()

	model := Model{
		updating:            false,
		textInput:           textInput,
		suggestions:         []Suggest{},
		filteredSuggestions: []Suggest{},
	}

	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case completionMsg:
		m.updating = false
		m.filteredSuggestions = msg

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.textInput.Value()

	if m.updating || m.prevText == m.textInput.Value() {
		return m, cmd
	}
	m.prevText = m.textInput.Value()

	m.updating = true
	return m, tea.Batch(cmd, m.updateCompletions())
}

type completionMsg []Suggest

func (m Model) updateCompletions() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(100 * time.Millisecond)
		search := strings.ToLower(m.textInput.Value())
		filtered := []Suggest{}
		for _, s := range m.suggestions {
			if strings.HasPrefix(strings.ToLower(s.Name), search) {
				filtered = append(filtered, s)
			}
		}

		return completionMsg(filtered)
	}
}

func (m Model) View() string {
	maxNameLen := 0
	maxDescLen := 0

	for _, s := range m.filteredSuggestions {
		if len(s.Name) > maxNameLen {
			maxNameLen = len(s.Name)
		}
		if len(s.Description) > maxDescLen {
			maxDescLen = len(s.Description)
		}
	}
	padding := lipgloss.NewStyle().PaddingLeft(m.textInput.Cursor() + len(m.textInput.Prompt)).Render("")
	nameStyle := lipgloss.
		NewStyle().
		PaddingLeft(1).
		Background(lipgloss.Color("8"))

	descStyle := nameStyle.
		Copy().
		Background(lipgloss.Color("9"))

	prompts := []string{m.textInput.View()}
	for _, s := range m.filteredSuggestions {
		name := nameStyle.PaddingRight(maxNameLen - len(s.Name) + 1).Render(s.Name)
		desc := descStyle.PaddingRight(maxDescLen - len(s.Description) + 1).Render(s.Description)
		line := lipgloss.JoinHorizontal(lipgloss.Bottom, padding, name, desc)
		prompts = append(prompts, line)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, prompts[:]...)
	return ret
}
