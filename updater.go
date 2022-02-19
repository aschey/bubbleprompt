package prompt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Order is important here, there's some strange freezing behavior
	// that happens if we update the text input before the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	// Scroll to bottom if the user typed something
	scrollToBottom := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSizeMsg(msg)

	case tea.KeyMsg:
		m.placeholderValue = ""
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		// Select next/previous list entry
		case tea.KeyUp, tea.KeyDown, tea.KeyTab:
			m.updateChosenListEntry(msg)

		case tea.KeyEnter:
			scrollToBottom = true
			cmds = m.submit(msg, cmds)

		case tea.KeyRunes, tea.KeyBackspace:
			scrollToBottom = true
			cmds = m.updateKeypress(msg, cmds)
		}

	case completionMsg:
		m.updating = false
		m.suggestions = msg

	case errMsg:
		m.err = msg
	}

	m.viewport.SetContent(m.render())
	if scrollToBottom {
		m.viewport.GotoBottom()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateWindowSizeMsg(msg tea.WindowSizeMsg) {
	if !m.ready {
		m.viewport = viewport.New(msg.Width, msg.Height-1)
		// TODO: register better bindings for these once the new input reader is merged
		m.viewport.KeyMap.Up = key.NewBinding(key.WithKeys("ctrl+a"))
		m.viewport.KeyMap.Down = key.NewBinding(key.WithKeys("ctrl+s"))
		m.ready = true
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
	}
}

func (m *Model) updateChosenListEntry(msg tea.KeyMsg) {
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
}

func (m *Model) submit(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	var curSuggestion *Suggestion
	if m.listPosition > -1 {
		curSuggestion = &m.suggestions[m.listPosition]
	}
	textValue := m.textInput.Value()

	// Reset all text and selection state
	m.textInput.SetValue("")
	m.typedText = ""
	m.listPosition = -1

	executorValue := m.executor(textValue, curSuggestion, m.suggestions)

	// Store the whole user input including the prompt state and the executor result
	// However note that we don't include all of textinput.View() because we don't want to include the cursor
	commandResult := lipgloss.JoinVertical(lipgloss.Left, m.textInput.Prompt+textValue, executorValue)
	m.previousCommands = append(m.previousCommands, commandResult)

	return append(cmds, m.updateCompletions())
}

func (m *Model) updateKeypress(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	m.typedText = m.textInput.Value()
	// Unselect selected item since user has changed the input
	m.listPosition = -1

	// If completer is already running or the text input hasn't changed, don't run the completer again
	if !m.updating && m.prevText != m.textInput.Value() {
		m.updating = true
		// Store last text ran against completer to compare against next time
		m.prevText = m.textInput.Value()
		cmds = append(cmds, m.updateCompletions())
	}

	return cmds
}
