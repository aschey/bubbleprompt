package prompt

import (
	"reflect"
	"strings"

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

	if m.executorModel != nil {
		var newCmd tea.Cmd
		executorModel := *m.executorModel
		executorModel, newCmd = executorModel.Update(msg)
		m.executorModel = &executorModel
		scrollToBottom = true

		// Check if the model sent the quit command
		// When this happens we just want to quit the executor, not the entire program
		// The only way to do this reliably without actually invoking the function is
		// to use reflection to check that the address is equal to tea.Quit's address
		if newCmd != nil && reflect.ValueOf(newCmd).Pointer() == reflect.ValueOf(tea.Quit).Pointer() {
			textValue := m.textInput.Value()
			executorValue := executorModel.View()

			// Store the whole user input including the prompt state and the executor result
			// However note that we don't include all of textinput.View() because we don't want to include the cursor
			commandResult := lipgloss.JoinVertical(lipgloss.Left, m.textInput.Prompt+textValue, executorValue)
			m.previousCommands = append(m.previousCommands, commandResult)
			m.executorModel = nil
			// Reset text after executor finished
			m.textInput.SetValue("")
			cmds = append(cmds, m.updateCompletions())
		} else {
			cmds = append(cmds, newCmd)
		}
	} else {
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
	// We'll reset the text input after the executor finished
	// so we can capture the final output
	m.typedText = ""
	m.listPosition = -1

	executorModel := m.executor(textValue, curSuggestion, m.suggestions)
	m.executorModel = &executorModel

	return append(cmds, executorModel.Init(), m.updateCompletions())
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

	prompts := m.previousCommands
	suggestionLength := len(m.suggestions)

	if m.executorModel != nil {
		// Executor is running, render next executor view
		// We're not going to render suggestions here, so set the length to 0 to apply the appropriate padding below the output
		suggestionLength = 0
		textView := m.textInput.Prompt + m.textInput.Value() + m.Placeholder.format(m.placeholderValue)
		prompts = append(prompts, textView)
		executorModel := *m.executorModel
		// Add a newline to ensure the text gets pushed up
		// this ensures the text doesn't jump if the completer takes a while to finish
		prompts = append(prompts, executorModel.View()+"\n")
	} else {
		textView := m.textInput.View() + m.Placeholder.format(m.placeholderValue)
		prompts = append(prompts, textView)

		// If an item is selected, parse out the text portion and apply formatting
		if m.listPosition > -1 {
			prompt := m.textInput.Prompt
			value := m.textInput.Value()
			formattedSuggestion := m.SelectedSuggestion.format(value)
			remainder := textView[len(prompt)+len(value):]
			textView = prompt + formattedSuggestion + remainder

		}
		for i, s := range m.suggestions {
			selected := i == m.listPosition
			name := m.Name.format(s.Name, maxNameLen, selected)
			description := m.Description.format(s.Description, maxDescLen, selected)

			line := lipgloss.JoinHorizontal(lipgloss.Bottom, padding, name, description)
			prompts = append(prompts, line)
		}

	}

	// Reserve height for prompts that were filtered out
	extraHeight := 5 - suggestionLength - 1
	if extraHeight > 0 {
		extraLines := strings.Repeat("\n", extraHeight)
		prompts = append(prompts, extraLines)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, prompts...)
	return ret
}
