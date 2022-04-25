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
	"github.com/charmbracelet/lipgloss"
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
	time.Sleep(100 * time.Millisecond)
	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), m.suggestions)
}

func executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "test"
	}), nil
}

func main() {
	placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	argStyle1 := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	argStyle2 := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	suggestions := []input.Suggestion[cmdMetadata]{
		{
			Text:        "first-option",
			Description: "test description",
			Metadata: commandinput.NewCmdMetadata([]commandinput.PositionalArg{
				{Placeholder: "<test1>", PlaceholderStyle: input.Text{Style: placeholderStyle}, ArgStyle: input.Text{Style: argStyle1}},
				{Placeholder: "<test2>", PlaceholderStyle: input.Text{Style: placeholderStyle}, ArgStyle: input.Text{Style: argStyle2}}}, commandinput.Placeholder{}),
		},
		{
			Text:        "second-option",
			Description: "test description2",
		},
		{
			Text:        "third-option",
			Description: "test description3",
		},
		{
			Text:        "fourth-option",
			Description: "test description4",
		},
		{
			Text:        "fifth-option",
			Description: "test description5",
			Metadata: commandinput.NewCmdMetadata([]commandinput.PositionalArg{
				{Placeholder: "<abc>", PlaceholderStyle: input.Text{Style: placeholderStyle}}}, commandinput.Placeholder{}),
		}}

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
