package bubbleprompt

import (
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
			{name: "test name", description: "test desc"},
			{name: "test name2", description: "test desc2"},
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
	padding := lipgloss.NewStyle().PaddingLeft(m.textInput.Cursor() + len(m.textInput.Prompt)).Render("")
	nameStyle := lipgloss.
		NewStyle().
		Background(lipgloss.Color("8"))
	descStyle := lipgloss.
		NewStyle().
		Background(lipgloss.Color("9"))

	prompts := []string{m.textInput.View()}
	for _, s := range m.suggestions {
		name := nameStyle.Render(s.name + " ")
		desc := descStyle.Render(s.description)
		line := lipgloss.JoinHorizontal(lipgloss.Left, padding, name, desc)
		prompts = append(prompts, line)
	}

	ret := lipgloss.JoinVertical(lipgloss.Top, prompts[:]...)
	return ret
}
