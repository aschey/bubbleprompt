package tutorial

import (
	prompt "github.com/aschey/bubbleprompt"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	// Update the counter every time the user submits something
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		m.numChoices++
	}
	return m, nil
}
