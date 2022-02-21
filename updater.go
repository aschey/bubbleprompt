package prompt

import (
	"reflect"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Check for exit signals before anything else
	// to reduce chance of program becoming frozen
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Order is important here, there's some strange freezing behavior
	// that happens if we update the text input before the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.completer, cmd = m.completer.Update(msg)
	cmds = append(cmds, cmd)

	// Scroll to bottom if the user typed something
	scrollToBottom := false

	switch m.modelState {
	case executing:
		cmds, scrollToBottom = m.updateExecuting(msg, cmds)
	case completing:
		cmds, scrollToBottom = m.updateCompleting(msg, cmds)
	}

	m.viewport.SetContent(m.render())
	if scrollToBottom {
		m.viewport.GotoBottom()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateExecuting(msg tea.Msg, cmds []tea.Cmd) ([]tea.Cmd, bool) {
	// Don't process text input while executor is running
	if m.textInput.Focused() {
		m.textInput.Blur()
	}

	executorModel, cmd := (*m.executorModel).Update(msg)
	m.executorModel = &executorModel

	// Check if the model sent the quit command
	// When this happens we just want to quit the executor, not the entire program
	// The only way to do this reliably without actually invoking the function is
	// to use reflection to check that the address is equal to tea.Quit's address
	if cmd != nil && reflect.ValueOf(cmd).Pointer() == reflect.ValueOf(tea.Quit).Pointer() {
		m.finalizeExecutor(executorModel)
		return cmds, true
	} else {
		return append(cmds, cmd), true
	}
}

func (m *Model) updateCompleting(msg tea.Msg, cmds []tea.Cmd) ([]tea.Cmd, bool) {
	scrollToBottom := false
	// Ensure text input is processing while executor is not running
	if !m.textInput.Focused() {
		cmds = append(cmds, m.textInput.Focus())
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSizeMsg(msg)

	case tea.KeyMsg:
		switch msg.Type {

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

	case errMsg:
		m.err = msg
	}
	return cmds, scrollToBottom
}

func (m *Model) finalizeExecutor(executorModel tea.Model) {
	textValue := m.textInput.Value()
	executorValue := executorModel.View()

	// Store the whole user input including the prompt state and the executor result
	// However note that we don't include all of textinput.View() because we don't want to include the cursor
	commandResult := lipgloss.JoinVertical(lipgloss.Left, m.textInput.Prompt+textValue, executorValue)
	m.previousCommands = append(m.previousCommands, commandResult)
	m.updateExecutor(nil)
	// Reset text after executor finished
	m.textInput.SetValue("")
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
	} else if (msg.Type == tea.KeyDown || msg.Type == tea.KeyTab) && m.listPosition < len(m.completer.suggestions)-1 {
		m.listPosition++
	} else {
		// -1 means no item selected
		m.listPosition = -1
	}

	if m.listPosition > -1 {
		// Set the input to the suggestion's selected text
		curSuggestion := m.completer.suggestions[m.listPosition]
		m.textInput.SetValue(curSuggestion.Name)
	} else {
		// If no selection, set the text back to the last thing the user typed
		m.textInput.SetValue(m.typedText)
	}

	// Move cursor to the end of the line
	m.textInput.SetCursor(len(m.textInput.Value()))
}

func (m *Model) updateExecutor(executor *tea.Model) {
	m.executorModel = executor
	if m.executorModel == nil {
		m.modelState = completing
	} else {
		m.modelState = executing
	}
}

func (m *Model) submit(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	var curSuggestion *Suggestion
	if m.listPosition > -1 {
		curSuggestion = &m.completer.suggestions[m.listPosition]
	}
	textValue := m.textInput.Value()
	// Reset all text and selection state
	// We'll reset the text input after the executor finished so we can capture the final output
	m.typedText = ""
	m.listPosition = -1

	executorModel := m.executor(textValue, curSuggestion, m.completer.suggestions)
	// Performance optimization: if this is a string model, we don't need to go through the whole update cycle
	// Just call the view method once and finalize the result
	// This makes the output a little cleaner if the completer function is slow
	if stringModel, ok := executorModel.(StringModel); ok {
		m.finalizeExecutor(stringModel)
	} else {
		m.updateExecutor(&executorModel)
		cmds = append(cmds, executorModel.Init())
	}

	return append(cmds, m.completer.updateCompletions(""))
}

func (m *Model) updateKeypress(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	m.typedText = m.textInput.Value()
	// Unselect selected item since user has changed the input
	m.listPosition = -1
	cmds = append(cmds, m.completer.updateCompletions(m.textInput.Value()))

	return cmds
}
