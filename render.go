package prompt

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) getPlaceholder() string {
	if !m.isSuggestionSelected() {
		return ""
	}
	placeholderText := m.completer.suggestions[m.listPosition].Placeholder

	if len(placeholderText) == 0 {
		return placeholderText
	}
	// Need to add an extra space between the placeholder and the value if the curor isn't at the end
	if m.textInput.Cursor() < len(m.textInput.Value()) {
		placeholderText = " " + placeholderText
	}
	return m.Formatters.Placeholder.format(placeholderText)
}

func (m Model) renderExecuting(lines []string) []string {
	textView := m.textInput.Prompt + m.textInput.Value() + m.getPlaceholder()
	lines = append(lines, textView)
	executorModel := *m.executorModel
	// Add a newline to ensure the text gets pushed up
	// this ensures the text doesn't jump if the completer takes a while to finish
	lines = append(lines, executorModel.View()+"\n")

	return lines
}

func (m Model) renderCompleting(lines []string) []string {
	// If an item is selected, parse out the text portion and apply formatting
	if m.isSuggestionSelected() {
		m.textInput.TextStyle = m.Formatters.SelectedSuggestion
	} else {
		m.textInput.TextStyle = lipgloss.NewStyle()
	}
	textView := m.textInput.View() + m.getPlaceholder()
	lines = append(lines, textView)

	// Calculate left offset for suggestions
	// Choosing a prompt via arrow keys or tab shouldn't change the prompt position
	// so we use the last typed cursor position instead of the current position
	paddingSize := len(m.textInput.Prompt) + m.lastTypedCursorPosition
	prompts := m.completer.suggestions.render(paddingSize, m.listPosition, m.Formatters)
	lines = append(lines, prompts...)

	return lines
}

func (m Model) render() string {
	lines := m.previousCommands
	suggestionLength := len(m.completer.suggestions)

	switch m.modelState {
	case executing:
		// Executor is running, render next executor view
		// We're not going to render suggestions here, so set the length to 0 to apply the appropriate padding below the output
		suggestionLength = 0
		lines = m.renderExecuting(lines)

	case completing:
		lines = m.renderCompleting(lines)
	}

	// Reserve height for prompts that were filtered out
	extraHeight := 5 - suggestionLength - 1
	if extraHeight > 0 {
		extraLines := strings.Repeat("\n", extraHeight)
		lines = append(lines, extraLines)
	}

	ret := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return ret
}
