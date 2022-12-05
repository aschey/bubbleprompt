package example

import (
	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	suggestions []suggestion.Suggestion[any]
	textInput   *simpleinput.Model[any]
}

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	return nil, nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	return nil, nil
}
