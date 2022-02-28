package prompt

import (
	"encoding/csv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) cursorView(v string, s lipgloss.Style) string {
	if m.textInput.Blink() {
		return s.Render(v)
	}
	return m.textInput.CursorStyle.Inline(true).Reverse(true).Render(v)
}

func (m Model) viewInput() string {
	args := []string{"testArg1", "testArg2"}
	argStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	textModel := m.textInput
	styleText := textModel.TextStyle.Render

	value := textModel.Value()

	pos := textModel.Cursor()
	v := styleText(value[:pos])

	if pos < len(value) {
		v += m.cursorView(string(value[pos]), textModel.TextStyle) // cursor and text under it
		v += styleText(value[pos+1:])                              // text after cursor
		if strings.HasPrefix(textModel.Placeholder, value) {
			v += textModel.PlaceholderStyle.Render(textModel.Placeholder[len(value):])
		}
	} else if pos < len(textModel.Placeholder) && strings.HasPrefix(textModel.Placeholder, value) {
		v += m.cursorView(string(textModel.Placeholder[pos]), m.textInput.PlaceholderStyle)
		v += textModel.PlaceholderStyle.Render(textModel.Placeholder[pos+1:])
	} else if len(args) == 0 {
		v += m.cursorView(" ", textModel.TextStyle)
	}

	if len(args) > 0 {
		numWords := 0
		r := csv.NewReader(strings.NewReader(textModel.Value()))
		r.Comma = ' '
		r.LazyQuotes = true
		record, _ := r.Read()
		for _, w := range record {
			if len(w) > 0 {
				numWords++
			}
		}
		if numWords == 0 {
			numWords = 1
		}

		argView := strings.Join(args[numWords-1:], " ")
		if !strings.HasSuffix(value, " ") && (pos < len(value) || value != textModel.Placeholder) {
			argView = " " + argView
		}

		if pos == len(value) {
			if len(argView) > 0 && !strings.HasPrefix(textModel.Placeholder, value) {
				v += m.cursorView(string(argView[0]), argStyle)
				v += argStyle.Render(argView[1:])
			} else {
				if !(pos < len(textModel.Placeholder) && strings.HasPrefix(textModel.Placeholder, textModel.Value())) {
					v += m.cursorView(" ", textModel.TextStyle)
				}
				v += argStyle.Render(argView)
			}
		} else {
			v += argStyle.Render(argView)
		}
	}

	return textModel.PromptStyle.Render(textModel.Prompt) + v
}

func (m Model) renderExecuting(lines []string) []string {
	textView := m.textInput.Prompt + m.textInput.Value()
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
	textView := m.viewInput()
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
