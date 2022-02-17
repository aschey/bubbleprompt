package prompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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

type Model struct {
	completer        Completer
	suggestions      []Suggestion
	textInput        textinput.Model
	Name             Text
	Description      Text
	Placeholder      Placeholder
	typedText        string
	prevText         string
	updating         bool
	listPosition     int
	placeholderValue string
	err              error
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
		Placeholder: Placeholder{
			ForegroundColor: "10",
			BackgroundColor: "",
		},
		suggestions:      completer(""),
		listPosition:     -1,
		placeholderValue: "",
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
	// Update text input if the user typed something
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.placeholderValue = ""
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		// Select next/previous list entry
		case tea.KeyUp, tea.KeyDown:
			if msg.Type == tea.KeyUp && m.listPosition > -1 {
				m.listPosition--
			} else if msg.Type == tea.KeyDown && m.listPosition < len(m.suggestions)-1 {
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
			return m, cmd

		case tea.KeyRunes:
			m.typedText = m.textInput.Value()

			if m.updating || m.prevText == m.textInput.Value() {
				return m, cmd
			}
			m.prevText = m.textInput.Value()

			m.updating = true
			return m, tea.Batch(cmd, m.updateCompletions())

		default:
			return m, cmd
		}

	case completionMsg:
		m.updating = false
		m.suggestions = msg
		return m, cmd

	case errMsg:
		m.err = msg
		return m, cmd

	default:
		return m, cmd
	}

}

type completionMsg []Suggestion

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
	padding := lipgloss.NewStyle().PaddingLeft(len(m.typedText) + len(m.textInput.Prompt)).Render("")

	textView := m.textInput.View() + m.Placeholder.format(m.placeholderValue)

	prompts := []string{textView}
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
