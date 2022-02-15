package prompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Suggest struct {
	Name        string
	Description string
}

type errMsg error

type Formatter func(name string, columnWidth int) string

type Completer func(input string) []Suggest

type Model struct {
	completer                  Completer
	suggestions                []Suggest
	textInput                  textinput.Model
	NameForegroundColor        string
	NameBackgroundColor        string
	NameFormatter              Formatter
	DescriptionForegroundColor string
	DescriptionBackgroundColor string
	DescriptionFormatter       Formatter
	prevText                   string
	updating                   bool
	err                        error
}

func New(completer Completer, opts ...Option) Model {
	textInput := textinput.New()
	textInput.Focus()

	model := Model{
		completer:                  completer,
		updating:                   false,
		textInput:                  textInput,
		NameForegroundColor:        "",
		NameBackgroundColor:        "14",
		DescriptionForegroundColor: "",
		DescriptionBackgroundColor: "37",
		suggestions:                completer(""),
	}

	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
}

func FilterHasPrefix(search string, suggestions []Suggest) []Suggest {
	lowerSearch := strings.ToLower(search)
	filtered := []Suggest{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(s.Name), lowerSearch) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

func (m Model) SetPrompt(prompt string) {
	m.textInput.Prompt = prompt
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
		m.suggestions = msg

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
		filtered := m.completer(m.textInput.Value())

		return completionMsg(filtered)
	}
}

func (m Model) View() string {
	maxNameLen := 0
	maxDescLen := 0

	for _, s := range m.suggestions {
		if len(s.Name) > maxNameLen {
			maxNameLen = len(s.Name)
		}
		if len(s.Description) > maxDescLen {
			maxDescLen = len(s.Description)
		}
	}
	padding := lipgloss.NewStyle().PaddingLeft(m.textInput.Cursor() + len(m.textInput.Prompt)).Render("")
	defaultStyle := lipgloss.
		NewStyle().
		PaddingLeft(1)

	prompts := []string{m.textInput.View()}
	for _, s := range m.suggestions {
		var name string
		if m.NameFormatter == nil {
			name = defaultStyle.
				Copy().
				Foreground(lipgloss.Color(m.NameForegroundColor)).
				Background(lipgloss.Color(m.NameBackgroundColor)).
				PaddingRight(maxNameLen - len(s.Name) + 1).
				Render(s.Name)
		} else {
			name = m.NameFormatter(s.Name, maxNameLen)
		}

		var desc string
		if m.DescriptionFormatter == nil {
			desc = defaultStyle.
				Copy().
				Foreground(lipgloss.Color(m.DescriptionForegroundColor)).
				Background(lipgloss.Color(m.DescriptionBackgroundColor)).
				PaddingRight(maxDescLen - len(s.Description) + 1).
				Render(s.Description)
		} else {
			desc = m.DescriptionFormatter(s.Description, maxDescLen)
		}

		line := lipgloss.JoinHorizontal(lipgloss.Bottom, padding, name, desc)
		prompts = append(prompts, line)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, prompts[:]...)
	return ret
}
