package prompt

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/internal"
	"github.com/charmbracelet/lipgloss"
)

func (m Model[T]) renderExecuting() string {
	executorModel := *m.executorModel
	// Add a newline to ensure the text gets pushed up
	// this ensures the text doesn't jump if the completer takes a while to finish
	textView := executorModel.View()
	return internal.AddNewlineIfMissing(textView)
}

func (m Model[T]) renderCompleting() string {
	// If an item is selected, parse out the text portion and apply formatting
	textView := internal.AddNewlineIfMissing(m.textInput.View(input.Interactive))

	// Calculate left offset for suggestions
	// Choosing a prompt via arrow keys or tab shouldn't change the prompt position
	// so we use the last typed cursor position instead of the current position
	paddingSize := len(m.textInput.Prompt()) + m.lastTypedCursorPosition
	prompts := m.completer.Render(paddingSize, m.formatters)
	textView += prompts

	return textView
}

func (m Model[T]) render() string {
	suggestionLength := len(m.completer.suggestions)
	if suggestionLength < 1 {
		// Always add at least one empty line
		suggestionLength = 1
	}
	lines := ""
	switch m.modelState {
	case executing:
		// Executor is running, render next executor view
		// We're not going to render suggestions here, so set the length to 0 to apply the appropriate padding below the output
		suggestionLength = 0
		lines = m.renderExecuting()

	case completing:
		lines = m.renderCompleting()
	}

	// Reserve height for prompts that were filtered out
	extraHeight := m.completer.maxSuggestions - suggestionLength
	if extraHeight > 0 {
		lines += strings.Repeat("\n", extraHeight)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, lines)
	return ret
}
