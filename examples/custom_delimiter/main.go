package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
)

type cmdMetadata = commandinput.CmdMetadata

type model struct {
	prompt prompt.Model[cmdMetadata]
}

type completerModel struct {
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[cmdMetadata]) []input.Suggestion[cmdMetadata] {
	if m.textInput.CommandCompleted() {
		suggestions := []input.Suggestion[cmdMetadata]{
			{Text: "abc"},
			{Text: "def"},
		}
		argText := ""
		if len(m.textInput.ParsedValue().Args.Value) > 0 {
			argText = m.textInput.ParsedValue().Args.Value[0].Value
		}
		return completers.FilterHasPrefix(argText, suggestions)
	}
	return completers.FilterHasPrefix(document.TextBeforeCursor(), m.suggestions)
}

func executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "test", nil
	}), nil
}

func main() {
	suggestions := []input.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test description"},
		{Text: "second-option", Description: "test description2"},
		{Text: "third-option", Description: "test description3"},
		{Text: "fourth-option", Description: "test description4"},
		{Text: "fifth-option", Description: "test description5"},
	}

	var textInput input.Input[cmdMetadata] = commandinput.New(commandinput.WithPrompt[cmdMetadata](">>> "),
		commandinput.WithDelimiterRegex[cmdMetadata](regexp.MustCompile(`[\s\.]+`)),
		commandinput.WithStringRegex[cmdMetadata](regexp.MustCompile(`[^\s\.]+`)))
	completerModel := completerModel{suggestions: suggestions, textInput: textInput.(*commandinput.Model[cmdMetadata])}

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
