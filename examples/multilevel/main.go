package main

import (
	"fmt"
	"os"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt prompt.Model
}

type completerModel struct {
	suggestions       []input.Suggestion
	textInput         *commandinput.Model
	filepathCompleter completers.FilePathCompleter
}

func (m model) Init() tea.Cmd {
	return m.prompt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.prompt.Update(msg)
	m.prompt = p
	return m, cmd
}

func (m model) View() string {
	return m.prompt.View()
}

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model) []input.Suggestion {
	if m.textInput.CommandCompleted() {
		filepath := ""
		parsed := m.textInput.ParsedValue()
		if len(parsed.Args.Value) > 0 {
			filepath = m.textInput.CurrentTokenBeforeCursor()
		}
		return m.filepathCompleter.Complete(filepath)
	}
	return completers.FilterHasPrefix(document.TextBeforeCursor(), m.suggestions)
}

func executor(input string, selected *input.Suggestion, suggestions []input.Suggestion) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "test"
	}), nil
}

func main() {
	suggestions := []input.Suggestion{
		{Text: "first-option", Description: "test description"},
		{Text: "second-option", Description: "test description2"},
		{Text: "third-option", Description: "test description3"},
		{Text: "fourth-option", Description: "test description4"},
		{Text: "fifth-option", Description: "test description5"},
	}

	textInput := commandinput.New()
	completerModel := completerModel{suggestions: suggestions, textInput: textInput}

	m := model{prompt: prompt.New(
		completerModel.completer,
		executor,
		textInput,
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
