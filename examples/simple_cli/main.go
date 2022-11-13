package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
)

type cmdMetadata struct {
	commandinput.CmdMetadata
	children []input.Suggestion[cmdMetadata]
}

func (c cmdMetadata) Children() []input.Suggestion[cmdMetadata] {
	return c.children
}

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
	parsed := m.textInput.ParsedValue()
	completed := m.textInput.CompletedArgsBeforeCursor()
	if len(completed) == 1 && parsed.Command.Value == "get" && parsed.Args.Value[0].Value == "weather" {
		flags := []commandinput.Flag{
			{
				Short:       "d",
				Long:        "days",
				Placeholder: m.textInput.NewFlagPlaceholder("<int>"),
				Description: "Forecast days",
			},
		}
		return m.textInput.FlagSuggestions(m.textInput.CurrentTokenBeforeCursor(), flags, nil), nil
	}
	return completer.GetRecursiveCompletions(m.textInput.Tokens(), m.textInput.CursorIndex(), m.suggestions), nil
}

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
	parsed := m.textInput.ParsedValue()
	if len(parsed.Args.Value) != 1 {
		return nil, fmt.Errorf("1 argument required")
	}
	return executor.NewStringModel("test"), nil
}

func main() {
	textInput := commandinput.New[cmdMetadata]()
	commandMetadata := commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<command>"))
	suggestions := []input.Suggestion[cmdMetadata]{
		{
			Text:        "get",
			Description: "retrieve things",
			Metadata: cmdMetadata{
				CmdMetadata: commandMetadata,
				children: []input.Suggestion[cmdMetadata]{
					{
						Text:        "secret",
						Description: "get the secret",
						Metadata: cmdMetadata{
							CmdMetadata: commandinput.CmdMetadata{Level: 1},
						},
					},
					{
						Text:        "weather",
						Description: "get the weather",
						Metadata: cmdMetadata{
							CmdMetadata: commandinput.CmdMetadata{Level: 1, ShowFlagPlaceholder: true},
						},
					},
				},
			},
		},
		{
			Text:        "set",
			Description: "update things",
			Metadata: cmdMetadata{
				CmdMetadata: commandMetadata,
				children: []input.Suggestion[cmdMetadata]{
					{
						Text:        "secret",
						Description: "update the secret",
						Metadata: cmdMetadata{
							CmdMetadata: commandinput.CmdMetadata{Level: 1},
						},
					},
				},
			},
		},
	}
	completerModel := completerModel{suggestions: suggestions, textInput: textInput}

	promptModel, err := prompt.New(
		completerModel.completer,
		completerModel.executor,
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
