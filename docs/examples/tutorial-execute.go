package tutorial

import (
	"fmt"
	"strconv"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/executor"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	// Get a list of all the tokens from the input
	tokens := m.textInput.WordTokenValues()
	if len(tokens) == 0 {
		// We didn't receive any input, which is invalid
		// Returning an error will output text will special error styling
		return nil, fmt.Errorf("No selection")
	}
	// The user entered a selection
	// Render their choice with styling applied
	return executor.NewStringModel(m.formatOutput(tokens[0])), nil
}

func (m model) formatOutput(choice string) string {
	return fmt.Sprintf("You picked: %s\nYou've entered %s submissions(s)\n\n",
		m.outputStyle.Render(choice),
		m.outputStyle.Render(strconv.FormatInt(m.numChoices, 10)))
}
