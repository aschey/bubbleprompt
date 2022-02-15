package bubbleprompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type suggest struct {
	name        string
	description string
}

type errMsg error

type Model struct {
	suggestions []suggest
	textInput   textinput.Model
	err         error
}

func New() Model {
	ti := textinput.New()
	ti.Prompt = ">>> "
	ti.Focus()
	return Model{
		textInput: ti,
		suggestions: []suggest{
			{name: "first option", description: "test desc"},
			{name: "second option", description: "test desc2"},
			{name: "third option", description: "test desc2"},
			{name: "fourth option", description: "test desc2"},
			{name: "fifth option", description: "test desc2"},
		}}
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

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.textInput.Value()
	return m, cmd
}

func (m Model) View() string {
	maxNameLen := 0
	maxDescLen := 0
	search := strings.ToLower(m.textInput.Value())
	filtered := []suggest{}
	for _, s := range m.suggestions {
		if strings.HasPrefix(strings.ToLower(s.name), search) {
			filtered = append(filtered, s)
		}
	}

	for _, s := range filtered {
		if len(s.name) > maxNameLen {
			maxNameLen = len(s.name)
		}
		if len(s.description) > maxDescLen {
			maxDescLen = len(s.description)
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
	for _, s := range filtered {
		name := nameStyle.PaddingRight(maxNameLen - len(s.name) + 1).Render(s.name)
		desc := descStyle.PaddingRight(maxDescLen - len(s.description) + 1).Render(s.description)
		line := lipgloss.JoinHorizontal(lipgloss.Bottom, padding, name, desc)
		prompts = append(prompts, line)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, prompts[:]...)
	return ret
}
