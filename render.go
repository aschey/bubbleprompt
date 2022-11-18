package prompt

import (
	"strings"

	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/internal"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func (m Model[T]) renderExecuting() string {
	executionManager := *m.executionManager
	// Add a newline to ensure the text gets pushed up
	// this ensures the text doesn't jump if the completer takes a while to finish
	textView := executionManager.View()
	return internal.AddNewlineIfMissing(textView)
}

func (m Model[T]) renderCompleting() string {
	// If an item is selected, parse out the text portion and apply formatting
	textView := internal.AddNewlineIfMissing(m.textInput.View(editor.Interactive))

	// Calculate left offset for suggestions
	// Choosing a prompt via arrow keys or tab shouldn't change the prompt position
	// so we use the last typed cursor position instead of the current position
	paddingSize := runewidth.StringWidth(m.textInput.Prompt()) + m.lastTypedCursorPosition
	prompts := m.suggestionManager.Render(paddingSize, m.formatters)
	textView += prompts

	return textView
}

func (m Model[T]) render() string {
	lines := ""
	contentHeight := 0
	switch m.modelState {
	case executing:
		// Executor is running, render next executor view
		lines = m.renderExecuting()

		// Add one line to account for the prompt + suggestions
		contentHeight = internal.CountNewlines(lines) + 1

	case completing:
		contentHeight = len(m.suggestionManager.suggestions)
		if contentHeight < 1 {
			// Always add at least one empty line
			contentHeight = 1
		}
		lines = m.renderCompleting()
	}

	// Reserve height for the max number of suggestions so the output height stays consistent
	extraHeight := m.suggestionManager.maxSuggestions - contentHeight
	if extraHeight > 0 {
		lines += strings.Repeat("\n", extraHeight)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, lines)
	return ret
}
