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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	if m.textInput.CommandCompleted() {
		suggestions := []input.Suggestion[cmdMetadata]{
			{Text: "abc"},
			{Text: "def"},
		}
		argText := ""
		if len(m.textInput.ParsedValue().Args.Value) > 0 {
			argText = m.textInput.ParsedValue().Args.Value[0].Value
		}
		return completers.FilterHasPrefix(argText, suggestions), nil
	}
	return completers.FilterHasPrefix(document.TextBeforeCursor(), m.suggestions), nil
}

func executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
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

	textInput := commandinput.New(commandinput.WithPrompt[cmdMetadata](">>> "),
		commandinput.WithDelimiterRegex[cmdMetadata](regexp.MustCompile(`[\s\.]+`)),
		commandinput.WithStringRegex[cmdMetadata](regexp.MustCompile(`[^\s\.]+`)))
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
