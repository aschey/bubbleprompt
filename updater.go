package prompt

import (
	"reflect"

	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model[I]) Update(msg tea.Msg) (Model[I], tea.Cmd) {
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
	m.renderer, cmd = m.renderer.Update(msg)
	cmds = append(cmds, cmd)

	prevText := m.textInput.Value()
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
		cmds, scrollToBottom = m.updateCompleting(msg, cmds, prevText)
	}

	cmd = m.finishUpdate(msg)
	cmds = append(cmds, cmd)

	m.renderer.SetContent(m.render())
	cmd = m.renderer.FinishUpdate()
	cmds = append(cmds, cmd)

	if scrollToBottom {
		m.renderer.GotoBottom(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model[I]) updateExecuting(msg tea.Msg, cmds []tea.Cmd) ([]tea.Cmd, bool) {
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

func (m *Model[I]) updateCompleting(msg tea.Msg, cmds []tea.Cmd, prevText string) ([]tea.Cmd, bool) {
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

		case tea.KeyBackspace, tea.KeyDelete, tea.KeyRunes, tea.KeySpace, tea.KeyLeft, tea.KeyRight:
			cmds = m.updateKeypress(msg, cmds, prevText)
		}

	case errMsg:
		m.err = msg
	}

	return cmds, scrollToBottom
}

func (m *Model[I]) selectSingle() {
	// Programatically select the suggestion if it's the only one and the input matches the suggestion
	if len(m.completer.suggestions) == 1 && m.textInput.ShouldSelectSuggestion(m.completer.suggestions[0]) {
		m.completer.selectSuggestion(m.completer.suggestions[0])
	}
}

func (m *Model[I]) finishUpdate(msg tea.Msg) tea.Cmd {
	suggestion := m.completer.getSelectedSuggestion()
	if suggestion == nil {
		// Nothing selected
		// Select the first suggestion if it matches
		m.selectSingle()

		cursor := m.textInput.Cursor()
		text := m.typedText
		// Get completion text before the cursor
		if cursor < len(text) {
			text = text[:cursor]
		}
		typedCompletionText := m.textInput.CompletionText(text)
		filteredSuggestions := completer.FilterHasPrefix(typedCompletionText, m.completer.suggestions)
		// Show placeholders for the first matching suggestion, but don't actually select it
		if len(filteredSuggestions) > 0 {
			suggestion = &filteredSuggestions[0]
		}
	}

	return m.textInput.OnUpdateFinish(msg, suggestion)
}

func (m *Model[I]) finalizeExecutor(executorModel *executorModel) {
	m.completer.unselectSuggestion()
	// Store the final executor view in the history
	// Need to store previous lines in a string instead of a []string in order to handle newlines from the tea.Model's View value properly
	// When executing a tea.Model standalone, the output must end in a newline and if we use a []string to track newlines, we'll get a double newline here
	m.renderer.AddOutput(executorModel.View())
	m.textInput.OnExecutorFinished()
	m.updateExecutor(nil)
}

func (m *Model[I]) updateWindowSizeMsg(msg tea.WindowSizeMsg) {
	if !m.ready {
		m.renderer.Initialize(msg)
		m.ready = true
	} else {
		m.renderer.SetSize(msg)
	}
}

func (m *Model[I]) updateChosenListEntry(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	if !m.completer.isSuggestionSelected() {
		// No suggestion currently suggested, store the last cursor position before selecting
		// so we can restore it later
		m.lastTypedCursorPosition = m.textInput.Cursor()
	}
	// Set the text back to the last thing the user typed in case the current suggestion changed the text length
	m.textInput.SetValue(m.typedText)
	// Make sure to set the cursor AFTER setting the value or it may get overwritten
	m.textInput.SetCursor(m.lastTypedCursorPosition)

	if msg.Type == tea.KeyUp {
		m.completer.previousSuggestion()
	} else {
		m.completer.nextSuggestion()
	}

	if m.completer.isSuggestionSelected() {
		// Set the input to the suggestion's selected text
		return nil
	} else {

		// Need to update completions since we changed the text and the cursor position
		return append(cmds, m.completer.updateCompletions(*m))
	}
}

func (m *Model[I]) updateExecutor(executor *executorModel) {
	m.executorModel = executor
	if m.executorModel == nil {
		m.modelState = completing
	} else {
		m.modelState = executing
	}
}

func (m *Model[I]) submit(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	innerExecutor, err := m.executor(m.textInput.Value())
	// Reset all text and selection state
	m.typedText = ""
	m.lastTypedCursorPosition = 0
	m.completer.unselectSuggestion()

	// Store the user input including the prompt state and the executor result
	// Pass in the static flag to signal to the text input to exclude interactive elements
	// such as placeholders and the cursor
	m.renderer.AddOutput(m.textInput.View(input.Static))
	m.textInput.SetValue("")

	executorModel := newExecutorModel(innerExecutor, m.Formatters.ErrorText, err)

	// Performance optimization: if this is a string model, we don't need to go through the whole update cycle
	// Just call the view method once and finalize the result
	// This makes the output a little cleaner if the completer function is slow
	if _, ok := innerExecutor.(executor.StringModel); ok {
		m.finalizeExecutor(executorModel)
	} else {
		m.updateExecutor(executorModel)
		cmds = append(cmds, executorModel.Init())
	}

	return append(cmds, m.completer.resetCompletions(*m))
}

func (m *Model[I]) updateKeypress(msg tea.KeyMsg, cmds []tea.Cmd, prevText string) []tea.Cmd {
	cmds = m.updatePosition(msg, cmds)
	if m.textInput.ShouldClearSuggestions(prevText, msg) {
		m.completer.clearSuggestions()
	} else if m.textInput.ShouldUnselectSuggestion(prevText, msg) {
		// Unselect selected item since user has changed the input
		m.completer.unselectSuggestion()
	}
	m.selectSingle()

	return cmds
}

func (m *Model[I]) updatePosition(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	m.lastTypedCursorPosition = m.textInput.Cursor()
	m.typedText = m.textInput.Value()
	cmds = append(cmds, m.completer.updateCompletions(*m))

	return cmds
}
