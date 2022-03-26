package prompt

import (
	"reflect"

	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

	cmd = m.textInput.OnUpdateStart(msg)
	cmds = append(cmds, cmd)

	m.completer, cmd = m.completer.Update(msg, m)
	cmds = append(cmds, cmd)

	// Scroll to bottom if the user typed something
	scrollToBottom := false

	switch m.modelState {
	case executing:
		cmds, scrollToBottom = m.updateExecuting(msg, cmds)
	case completing:
		cmds, scrollToBottom = m.updateCompleting(msg, cmds)
	}

	cmd = m.finishUpdate(msg)
	cmds = append(cmds, cmd)

	m.viewport.SetContent(m.render())
	if scrollToBottom {
		m.viewport.GotoBottom()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateExecuting(msg tea.Msg, cmds []tea.Cmd) ([]tea.Cmd, bool) {
	executorModel, cmd := (*m.executorModel).Update(msg)
	m.executorModel = &executorModel

	// Check if the model sent the quit command
	// When this happens we just want to quit the executor, not the entire program
	// The only way to do this reliably without actually invoking the function is
	// to use reflection to check that the address is equal to tea.Quit's address
	if cmd != nil && reflect.ValueOf(cmd).Pointer() == reflect.ValueOf(tea.Quit).Pointer() {
		m.finalizeExecutor(m.executorModel)
		// Re-focus input when finished
		return append(cmds, m.textInput.Focus()), true
	} else {
		// Don't process text input while executor is running
		if m.textInput.Focused() {
			m.textInput.Blur()
		}
		return append(cmds, cmd), true
	}
}

func (m *Model) updateCompleting(msg tea.Msg, cmds []tea.Cmd) ([]tea.Cmd, bool) {
	scrollToBottom := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSizeMsg(msg)

	case tea.KeyMsg:
		scrollToBottom = true
		switch msg.Type {

		// Select next/previous list entry
		case tea.KeyUp, tea.KeyDown, tea.KeyTab:
			cmds = m.updateChosenListEntry(msg, cmds)

		case tea.KeyEnter:
			cmds = m.submit(msg, cmds)

		case tea.KeyRunes, tea.KeyBackspace, tea.KeyDelete:
			cmds = m.updateKeypress(msg, cmds)

		case tea.KeyLeft, tea.KeyRight:
			cmds = m.updatePosition(msg, cmds)
		}

	case errMsg:
		m.err = msg
	}

	return cmds, scrollToBottom
}

func (m *Model) selectSingle(text string) {
	// Programatically select the suggestion if it's the only one and the text matches the suggestion
	completionText := m.textInput.CompletionText(text)
	if len(m.completer.suggestions) == 1 && completionText == m.completer.suggestions[0].Text {
		m.completer.selectedKey = m.completer.suggestions[0].Key()
	}
}

func (m *Model) finishUpdate(msg tea.Msg) tea.Cmd {
	suggestion := m.completer.getSelectedSuggestion()
	if suggestion == nil {
		// Nothing selected
		// Select the first suggestion if it matches
		m.selectSingle(m.textInput.Value()[:m.textInput.Cursor()])

		typedCompletionText := m.textInput.CompletionText(m.typedText[:m.textInput.Cursor()])
		filteredSuggestions := completer.FilterHasPrefix(typedCompletionText, m.completer.suggestions)
		// Show placeholders for the first matching suggestion, but don't actually select it
		if len(filteredSuggestions) > 0 {
			suggestion = &filteredSuggestions[0]
		}
	}
	return m.textInput.OnUpdateFinish(msg, suggestion)
}

func (m *Model) finalizeExecutor(executorModel *executorModel) {
	m.completer.unselectSuggestion()
	// Store the final executor view in the history
	// Need to store previous lines in a string instead of a []string in order to handle newlines from the tea.Model's View value properly
	// When executing a tea.Model standalone, the output must end in a newline and if we use a []string to track newlines, we'll get a double newline here
	m.previousCommands += executorModel.View()
	m.updateExecutor(nil, nil)
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

func (m *Model) updateChosenListEntry(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	if !m.completer.isSuggestionSelected() {
		// No suggestion currently suggested, store the last cursor position before selecting
		// so we can restore it later
		m.lastTypedCursorPosition = m.textInput.Cursor()
	}

	if msg.Type == tea.KeyUp {
		m.completer.previousSuggestion()
	} else {
		m.completer.nextSuggestion()
	}

	if m.completer.isSuggestionSelected() {
		// Set the input to the suggestion's selected text
		curSuggestion := m.completer.getSelectedSuggestion()
		m.textInput.OnSuggestionChanged(*curSuggestion)

		return nil
	} else {
		// If no selection, set the text back to the last thing the user typed
		m.textInput.SetValue(m.typedText)
		m.textInput.SetCursor(m.lastTypedCursorPosition)
		// Need to update completions since we changed the text and the cursor position
		return append(cmds, m.completer.updateCompletions(*m))
	}
}

func (m *Model) updateExecutor(executor *executorModel, err error) {
	m.executorModel = executor
	if m.executorModel == nil {
		m.modelState = completing
	} else {
		m.modelState = executing
	}
}

func (m *Model) submit(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	curSuggestion := m.completer.getSelectedSuggestion()
	textValue := m.textInput.Value()
	innerExecutor, err := m.executor(textValue, curSuggestion, m.completer.suggestions)
	// Reset all text and selection state
	m.typedText = ""
	m.lastTypedCursorPosition = 0
	m.completer.unselectSuggestion()

	// Store the whole user input including the prompt state and the executor result
	// However note that we don't include all of textInput.View() because we don't want to include the cursor
	m.previousCommands += m.textInput.Prompt() + textValue + "\n"
	m.textInput.SetValue("")

	executorModel := newExecutorModel(innerExecutor, m.Formatters.ErrorText, err)

	// Performance optimization: if this is a string model, we don't need to go through the whole update cycle
	// Just call the view method once and finalize the result
	// This makes the output a little cleaner if the completer function is slow
	if _, ok := innerExecutor.(executor.StringModel); ok {
		m.finalizeExecutor(executorModel)
	} else {
		m.updateExecutor(executorModel, err)
		cmds = append(cmds, executorModel.Init())
	}

	return append(cmds, m.completer.resetCompletions(*m))
}

func (m *Model) updateKeypress(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	cmds = m.updatePosition(msg, cmds)

	if !m.textInput.IsDelimiter(msg.String()) {
		// Unselect selected item since user has changed the input
		m.completer.unselectSuggestion()
	}
	m.selectSingle(m.textInput.Value())

	return cmds
}

func (m *Model) updatePosition(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	m.lastTypedCursorPosition = m.textInput.Cursor()
	m.typedText = m.textInput.Value()
	cmds = append(cmds, m.completer.updateCompletions(*m))

	return cmds
}
