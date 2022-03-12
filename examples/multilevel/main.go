package main

import (
	"fmt"
	"os"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt prompt.Model
}

type completerModel struct {
	suggestions       []prompt.Suggestion
	filepathCompleter prompt.FilePathCompleter
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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model) prompt.Suggestions {
	if promptModel.CommandCompleted() {
		filepath := ""
		if len(document.ParsedInput.Args.Value) > 0 {
			filepath = document.ParsedInput.Args.Value[0].Value
		}
		return m.filepathCompleter.Complete(filepath)
	}
	return prompt.FilterHasPrefix(document.TextBeforeCursor(), m.suggestions)
}

func executor(input string, selected *prompt.Suggestion, suggestions prompt.Suggestions) tea.Model {
	return prompt.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "test"
	})
}

func main() {
	suggestions := []prompt.Suggestion{
		{Text: "first-option", Description: "test description"},
		{Text: "second-option", Description: "test description2"},
		{Text: "third-option", Description: "test description3"},
		{Text: "fourth-option", Description: "test description4"},
		{Text: "fifth-option", Description: "test description5"},
	}

	completerModel := completerModel{suggestions: suggestions}

	m := model{prompt: prompt.New(
		completerModel.completer,
		executor,
		prompt.WithPrompt(">>> "),
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
