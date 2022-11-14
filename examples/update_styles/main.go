package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cmdMetadata struct {
	commandinput.CmdMetadata
	children []input.Suggestion[cmdMetadata]
}

func (c cmdMetadata) Children() []input.Suggestion[cmdMetadata] {
	return c.children
}

type completerModel struct {
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
}

func (m completerModel) completer(promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	return completer.GetRecursiveCompletions(m.textInput.Tokens(), m.textInput.CursorIndex(), m.suggestions), nil
}

func (m *completerModel) executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
	parsed := m.textInput.ParsedValue()
	if parsed.Command.Value == "cursor-mode" {
		switch parsed.Args.Value[0].Value {
		case "blink":
			return executor.NewCmdModel("blinking cursor", m.textInput.SetCursorMode(textinput.CursorBlink)), nil
		case "static":
			return executor.NewCmdModel("static cursor", m.textInput.SetCursorMode(textinput.CursorStatic)), nil
		case "hide":
			return executor.NewCmdModel("blinking cursor", m.textInput.SetCursorMode(textinput.CursorHide)), nil
		}

	}
	return nil, fmt.Errorf("Invalid input")
}

func main() {
	textInput := commandinput.New[cmdMetadata]()
	secretArgs := textInput.NewPositionalArgs("<secret value>")
	secretArgs[0].ArgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	commandMetadata := commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<command>"))

	suggestions := []input.Suggestion[cmdMetadata]{
		{
			Text:        "cursor-mode",
			Description: "set the cursor mode",
			Metadata: cmdMetadata{
				CmdMetadata: commandMetadata,
				children: []input.Suggestion[cmdMetadata]{
					{
						Text:        "blink",
						Description: "blinking cursor",
						Metadata: cmdMetadata{
							CmdMetadata: commandinput.CmdMetadata{
								Level: 1,
							},
						},
					},
					{
						Text:        "static",
						Description: "normal cursor",
						Metadata: cmdMetadata{
							CmdMetadata: commandinput.CmdMetadata{
								Level: 1,
							},
						},
					},
					{
						Text:        "hide",
						Description: "no cursor",
						Metadata: cmdMetadata{
							CmdMetadata: commandinput.CmdMetadata{
								Level: 1,
							},
						},
					},
				},
			},
		},
	}
	completerModel := completerModel{
		suggestions: suggestions,
		textInput:   textInput,
	}

	promptModel, err := prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	)
	if err != nil {
		panic(err)
	}

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
