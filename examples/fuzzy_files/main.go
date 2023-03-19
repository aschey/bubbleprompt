package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput   *simpleinput.Model[any]
	outputStyle lipgloss.Style
	filterer    completer.PathCompleter[any]
}

func (m model) Complete(promptModel prompt.Model[any]) ([]suggestion.Suggestion[any], error) {
	return m.filterer.Complete(m.textInput.CurrentTokenBeforeCursor()), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	return executor.NewStringModel(m.formatOutput(m.textInput.Value())), nil
}

func (m model) formatOutput(choice string) string {
	return fmt.Sprintf("You picked: %s\n\n",
		m.outputStyle.Render(choice),
	)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	return m, nil
}

func main() {
	textInput := simpleinput.New[any]()

	model := model{
		textInput:   textInput,
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
		filterer:    completer.PathCompleter[any]{Filterer: completer.NewFuzzyFilter[any]()},
	}

	promptModel := prompt.New[any](model, textInput)

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Search for files or directories"))
	fmt.Println()

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
