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

type Formatter func(name string, columnWidth int, selected bool) string

type Completer func(input string) []Suggest

type Text struct {
	ForegroundColor         string
	SelectedForegroundColor string
	BackgroundColor         string
	SelectedBackgroundColor string
	Formatter               Formatter
}

func (t Text) format(text string, maxLen int, selected bool) string {
	defaultStyle := lipgloss.
		NewStyle().
		PaddingLeft(1)

	foregroundColor := t.ForegroundColor
	backgroundColor := t.BackgroundColor
	if selected {
		foregroundColor = t.SelectedForegroundColor
		backgroundColor = t.SelectedBackgroundColor
	}
	var formattedText string
	if t.Formatter == nil {
		formattedText = defaultStyle.
			Copy().
			Foreground(lipgloss.Color(foregroundColor)).
			Background(lipgloss.Color(backgroundColor)).
			PaddingRight(maxLen - len(text) + 1).
			Render(text)
	} else {
		formattedText = t.Formatter(text, maxLen, selected)
	}

	return formattedText
}

type Model struct {
	completer    Completer
	suggestions  []Suggest
	textInput    textinput.Model
	Name         Text
	Description  Text
	prevText     string
	updating     bool
	listPosition int
	err          error
}

func New(completer Completer, opts ...Option) Model {
	textInput := textinput.New()
	textInput.Focus()

	model := Model{
		completer: completer,
		updating:  false,
		textInput: textInput,
		Name: Text{
			ForegroundColor:         "",
			SelectedForegroundColor: "240",
			BackgroundColor:         "14",
			SelectedBackgroundColor: "14",
		},
		Description: Text{
			ForegroundColor:         "",
			SelectedForegroundColor: "240",
			BackgroundColor:         "37",
			SelectedBackgroundColor: "37",
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
		case tea.KeyDown:
			if m.listPosition < len(m.suggestions)-1 {
				m.listPosition++
			} else {
				m.listPosition = -1
			}

		case tea.KeyUp:
			if m.listPosition > -1 {
				m.listPosition--
			}

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

	prompts := []string{m.textInput.View()}
	for i, s := range m.suggestions {
		selected := i == m.listPosition
		name := m.Name.format(s.Name, maxNameLen, selected)
		description := m.Description.format(s.Description, maxDescLen, selected)

		line := lipgloss.JoinHorizontal(lipgloss.Bottom, padding, name, description)
		prompts = append(prompts, line)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, prompts[:]...)
	return ret
}
