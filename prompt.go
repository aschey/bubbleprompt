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
	completer        Completer
	executor         Executor
	suggestions      []Suggestion
	textInput        textinput.Model
	viewport         viewport.Model
	history          []string
	Name             Text
	Description      Text
	Placeholder      Placeholder
	typedText        string
	prevText         string
	updating         bool
	listPosition     int
	placeholderValue string
	ready            bool
	err              error
}

func New(completer Completer, executor Executor, opts ...Option) Model {
	textInput := textinput.New()
	textInput.Focus()

	model := Model{
		completer: completer,
		executor:  executor,
		textInput: textInput,
		Name: Text{
			SelectedForegroundColor: "240",
			BackgroundColor:         "14",
			SelectedBackgroundColor: "14",
		},
		Description: Text{
			SelectedForegroundColor: "240",
			BackgroundColor:         "37",
			SelectedBackgroundColor: "37",
		},
		Placeholder: Placeholder{
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	// Update text input if the user typed something
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

	case tea.KeyMsg:
		m.placeholderValue = ""
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		// Select next/previous list entry
		case tea.KeyUp, tea.KeyDown, tea.KeyTab:
			if msg.Type == tea.KeyUp && m.listPosition > -1 {
				m.listPosition--
			} else if (msg.Type == tea.KeyDown || msg.Type == tea.KeyTab) && m.listPosition < len(m.suggestions)-1 {
				m.listPosition++
			} else {
				// -1 means no item selected
				m.listPosition = -1
			}

			if m.listPosition > -1 {
				// Set the input to the suggestion's selected text
				curSuggestion := m.suggestions[m.listPosition]
				m.placeholderValue = curSuggestion.Placeholder
				m.textInput.SetValue(curSuggestion.Name)
			} else {
				// If no selection, set the text back to the last thing the user typed
				m.textInput.SetValue(m.typedText)
			}

			// Move cursor to the end of the line
			m.textInput.SetCursor(len(m.textInput.Value()))

		case tea.KeyEnter:
			var curSuggestion *Suggestion
			if m.listPosition > -1 {
				curSuggestion = &m.suggestions[m.listPosition]
			}
			textValue := m.textInput.Value()
			m.textInput.SetValue("")
			executorValue := m.executor(textValue, curSuggestion, m.suggestions)

			ret := lipgloss.JoinVertical(lipgloss.Left, m.textInput.Prompt+textValue, executorValue)
			m.history = append(m.history, ret)
			cmds = append(cmds, m.updateCompletions())

		case tea.KeyRunes, tea.KeyBackspace:
			m.typedText = m.textInput.Value()

			if m.updating || m.prevText == m.textInput.Value() {
				return m, cmd
			}
			m.prevText = m.textInput.Value()
			m.updating = true
			cmds = append(cmds, m.updateCompletions())
			m.viewport.LineDown(2)
		}

	case completionMsg:
		m.updating = false
		m.suggestions = msg

	case errMsg:
		m.err = msg

	}

	m.viewport.SetContent(m.render())
	m.viewport.GotoBottom()
	return m, tea.Batch(cmds...)

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

	for _, s := range m.suggestions {
		if len(s.Name) > maxNameLen {
			maxNameLen = len(s.Name)
		}
		if len(s.Description) > maxDescLen {
			maxDescLen = len(s.Description)
		}
	}
	padding := lipgloss.NewStyle().PaddingLeft(len(m.typedText) + len(m.textInput.Prompt)).Render("")

	textView := m.textInput.View() + m.Placeholder.format(m.placeholderValue)

	prompts := append(m.history, textView)

	for i, s := range m.suggestions {
		selected := i == m.listPosition
		name := m.Name.format(s.Name, maxNameLen, selected)
		description := m.Description.format(s.Description, maxDescLen, selected)

		line := lipgloss.JoinHorizontal(lipgloss.Bottom, padding, name, description)
		prompts = append(prompts, line)
	}
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
