package prompt

import (
	"reflect"

	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		// Ctrl+C should always shutdown the whole program regardless
		// of what the executor is doing
		case tea.KeyCtrlC:
			shutdown = true
			return m, tea.Quit
		}
	case rendererMsg:
		// No need to switch renderers if they're the same type
		if reflect.TypeOf(m.renderer) != reflect.TypeOf(msg.renderer) {
			currentHistory := m.renderer.GetHistory()

			m.renderer = msg.renderer
			m.renderer.Initialize(m.size)

			if msg.retainHistory {
				cmds = append(cmds, m.renderer.SetHistory(currentHistory))
			}
		}
	case suggestion.CompleteMsg:
		cmds = append(cmds, func() tea.Msg {
			filtered, err := m.inputHandler.Complete(m)
			return suggestion.SuggestionMsg[T]{Suggestions: filtered, Err: err}
		})
	}

	m.inputHandler, cmd = m.inputHandler.Update(msg)
	cmds = append(cmds, cmd)

	// Order is important here
	m.renderer, cmd = m.renderer.Update(msg)
	cmds = append(cmds, cmd)

	prevText := m.textInput.Runes()
	cmd = m.textInput.OnUpdateStart(msg)
	cmds = append(cmds, cmd)

	if m.suggestionManager.ShouldChangeListPosition(msg) {
		m.saveCurrentInput()
	}

	cmds = append(cmds, m.suggestionManager.Update(msg))

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

	m.renderer.SetInput(m.renderInput())
	m.renderer.SetBody(m.renderBody())
	cmd = m.renderer.FinishUpdate()
	cmds = append(cmds, cmd)

	if scrollToBottom {
		m.renderer.GotoBottom(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model[T]) updateExecuting(msg tea.Msg, cmds []tea.Cmd) ([]tea.Cmd, bool) {
	executorManager, cmd := (*m.executionManager).Update(msg)
	m.executionManager = &executorManager

	switch msg.(type) {
	// Check if the model sent the quit command
	// When this happens we just want to quit the executor, not the entire program
	case quitAttempted:
		cmds = append(cmds, m.finalizeExecutor(m.executionManager))
		// Re-focus input when finished
		return append(cmds, m.textInput.Focus()), true
	default:
		if m.textInput.Focused() {
			m.textInput.Blur()
		}
	}
	return append(cmds, cmd), true
}

func (m *Model[T]) updateCompleting(
	msg tea.Msg,
	cmds []tea.Cmd,
	prevRunes []rune,
) ([]tea.Cmd, bool) {
	scrollToBottom := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSizeMsg(msg)
	case tea.KeyMsg:
		scrollToBottom = true
		switch msg.Type {
		case tea.KeyEscape:
			// Escape should only shutdown the program in completer mode, otherwise this could interfere
			// with the executor model
			shutdown = true
			return append(cmds, tea.Quit), scrollToBottom

		case tea.KeyEnter:
			cmds = m.submit(msg, cmds)

		case tea.KeyBackspace, tea.KeyDelete, tea.KeyRunes, tea.KeySpace, tea.KeyLeft, tea.KeyRight:
			cmds = m.updateKeypress(msg, cmds, prevRunes)
		}

	case errMsg:
		m.err = msg
	}

	return cmds, scrollToBottom
}

func (m *Model[T]) selectSingle() {
	// Programatically select the suggestion if it's the only one and the input matches the suggestion
	suggestions := m.suggestionManager.Suggestions()
	if len(suggestions) > 0 {
		firstSuggestion := suggestions[0]
		// Nothing selected
		// Select the first suggestion if it matches
		if m.suggestionManager.SelectedSuggestion() == nil && len(suggestions) == 1 &&
			m.textInput.ShouldSelectSuggestion(firstSuggestion) {
			m.suggestionManager.SelectSuggestion(firstSuggestion)
		}
	}
}

func (m *Model[T]) finishUpdate(msg tea.Msg) tea.Cmd {
	m.selectSingle()
	suggestion := m.suggestionManager.SelectedSuggestion()
	isSelected := suggestion != nil
	if !isSelected {
		// Show placeholders for the first matching suggestion, but don't actually select it
		if len(m.suggestionManager.Suggestions()) > 0 {
			suggestion = &m.suggestionManager.Suggestions()[0]
		}
	}

	return m.textInput.OnUpdateFinish(msg, suggestion, isSelected)
}

func (m *Model[T]) finalizeExecutor(executorManager *executionManager) tea.Cmd {
	m.suggestionManager.UnselectSuggestion()
	// Store the final executor view in the history
	// Need to store previous lines in a string instead of a []string in order
	// to handle newlines from the tea.Model's View value properly
	// When executing a tea.Model standalone, the output must end in a newline and
	// if we use a []string to track newlines, we'll get a double newline here
	m.renderer.AddHistory(executorManager.View())
	m.textInput.OnExecutorFinished()
	m.updateExecutor(nil)
	return func() tea.Msg { return ExecutorFinishedMsg(executorManager.inner) }
}

func (m *Model[T]) updateWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.size = msg
	if !m.ready {
		m.renderer.Initialize(msg)
		m.ready = true
	} else {
		m.renderer.SetSize(msg)
	}
}

func (m *Model[T]) saveCurrentInput() {
	if !m.suggestionManager.IsSuggestionSelected() {
		// No suggestion currently suggested, store the last cursor position before selecting
		// so we can restore it later
		m.lastTypedCursorPosition = m.textInput.CursorOffset()
	}
	// Set the text back to the last thing the user typed in case the current suggestion changed the text length
	m.textInput.SetValue(string(m.typedRunes))
	// Make sure to set the cursor AFTER setting the value or it may get overwritten
	m.textInput.SetCursor(m.lastTypedCursorPosition)
}

func (m *Model[T]) updateExecutor(executor *executionManager) {
	m.executionManager = executor
	if m.executionManager == nil {
		m.modelState = completing
	} else {
		m.modelState = executing
	}
}

func (m *Model[T]) submit(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	innerExecutor, err := m.inputHandler.Execute(m.textInput.Value(), m)
	if innerExecutor == nil {
		// No executor returned, default to empty model to prevent nil reference errors
		innerExecutor = executor.NewStringModel("")
	}
	// Reset all text and selection state
	m.typedRunes = []rune("")
	m.lastTypedCursorPosition = 0
	m.suggestionManager.UnselectSuggestion()

	// Store the user input including the prompt state and the executor result
	// Pass in the static flag to signal to the text input to exclude interactive elements
	// such as placeholders and the cursor
	m.renderer.AddHistory(m.textInput.View(input.Static))
	m.textInput.ResetValue()

	executorManager := newExecutorManager(innerExecutor, m.suggestionManager.Formatters().ErrorText, err)

	// Performance optimization: if this is a string model, we don't need to go through the whole update cycle
	// Just call the view method once and finalize the result
	// This makes the output a little cleaner if the completer function is slow
	if _, ok := innerExecutor.(executor.StringModel); ok {
		cmds = append(cmds, m.finalizeExecutor(executorManager))
	} else {
		m.updateExecutor(executorManager)
		// Need to explicitly notify the child model of the current window size.
		// Since the bubbletea event loop is already running, this won't happen automatically.
		cmds = append(cmds, tea.Sequence(executorManager.Init(), func() tea.Msg { return m.size }))
	}

	return append(cmds, m.suggestionManager.ResetSuggestions())
}

func (m *Model[T]) updateKeypress(msg tea.KeyMsg, cmds []tea.Cmd, prevRunes []rune) []tea.Cmd {
	cmds = m.updatePosition(msg, cmds)
	if m.textInput.ShouldClearSuggestions(prevRunes, msg) {
		m.suggestionManager.ClearSuggestions()
	} else if m.textInput.ShouldUnselectSuggestion(prevRunes, msg) {
		// Unselect selected item since user has changed the input
		m.suggestionManager.UnselectSuggestion()
	}
	m.selectSingle()

	return cmds
}

func (m *Model[T]) updatePosition(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	m.lastTypedCursorPosition = m.textInput.CursorOffset()
	m.typedRunes = m.textInput.Runes()
	cmds = append(cmds, m.suggestionManager.UpdateSuggestions())

	return cmds
}
