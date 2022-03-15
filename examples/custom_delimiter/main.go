package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt prompt.Model
}

type completerModel struct {
	suggestions input.Suggestions
	textInput   *commandinput.Model
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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model) input.Suggestions {
	if m.textInput.CommandCompleted() {
		suggestions := input.Suggestions{
			{Text: "abc"},
			{Text: "def"},
		}
		argText := ""
		if len(m.textInput.ParsedValue().Args.Value) > 0 {
			argText = m.textInput.ParsedValue().Args.Value[0].Value
		}
		return prompt.FilterHasPrefix(argText, suggestions)
	}
	return prompt.FilterHasPrefix(document.TextBeforeCursor(), m.suggestions)
}

func executor(input string, selected *input.Suggestion, suggestions input.Suggestions) tea.Model {
	return prompt.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "test"
	})
}

func main() {
	suggestions := input.Suggestions{
		{Text: "first-option", Description: "test description"},
		{Text: "second-option", Description: "test description2"},
		{Text: "third-option", Description: "test description3"},
		{Text: "fourth-option", Description: "test description4"},
		{Text: "fifth-option", Description: "test description5"},
	}

	textInput := commandinput.New(commandinput.WithPrompt(">>> "),
		commandinput.WithDelimiterRegex(regexp.MustCompile(`[\s\.]+`)),
		commandinput.WithStringRegex(regexp.MustCompile(`[^\s\.]+`)))
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
