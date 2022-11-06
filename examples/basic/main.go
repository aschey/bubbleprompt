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
	promptModel prompt.Model[cmdMetadata]
}

type completerModel struct {
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
}

func (m model) Init() tea.Cmd {
	return m.promptModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.promptModel.Update(msg)
	m.promptModel = p
	return m, cmd
}

func (m model) View() string {
	return m.promptModel.View()
}

func (m completerModel) completer(promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	time.Sleep(100 * time.Millisecond)
	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), m.suggestions), nil
}

func executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "result is " + input, nil
	}), nil
}

func main() {
	suggestions := []input.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test desc", Metadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("[test placeholder]")}}},
		{Text: "second-option", Description: "test desc2"},
		{Text: "third-option", Description: "test desc3"},
		{Text: "fourth-option", Description: "test desc4"},
		{Text: "fifth-option", Description: "test desc5"},
		{Text: "sixth-option", Description: "test desc6"},
		{Text: "seventh-option", Description: "test desc7"}}

	textInput := commandinput.New[cmdMetadata]()
	completerModel := completerModel{suggestions: suggestions, textInput: textInput}

	promptModel, err := prompt.New(
		completerModel.completer,
		executor,
		textInput,
	)
	if err != nil {
		panic(err)
	}
	m := model{promptModel}

	if _, err := tea.NewProgram(m, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
