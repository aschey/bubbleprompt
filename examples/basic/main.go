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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	time.Sleep(100 * time.Millisecond)
	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), m.suggestions), nil
}

func executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "result is " + input, nil
	}), nil
}

func main() {
	suggestions := []input.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test desc", Metadata: commandinput.NewCmdMetadata([]commandinput.PositionalArg{commandinput.NewPositionalArg("[test placeholder]")}, commandinput.Placeholder{})},
		{Text: "second-option", Description: "test desc2"},
		{Text: "third-option", Description: "test desc3"},
		{Text: "fourth-option", Description: "test desc4"},
		{Text: "fifth-option", Description: "test desc5"},
		{Text: "sixth-option", Description: "test desc6"},
		{Text: "seventh-option", Description: "test desc7"}}

	var textInput input.Input[cmdMetadata] = commandinput.New[cmdMetadata]()
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
