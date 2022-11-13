package main

import (
	"fmt"
	"os"
	"strconv"

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

type completerModel struct {
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
	secret      string
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

func (m *completerModel) executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
	parsed := m.textInput.ParsedValue()
	args := parsed.Args.Value
	flags := parsed.Flags.Value
	if len(args) == 0 {
		return nil, fmt.Errorf("1 argument required")
	}
	arg := args[0].Value
	switch parsed.Command.Value {
	case "get":
		switch arg {
		case "weather":
			days := int64(1)
			if len(flags) > 0 {
				flag := flags[0]
				if flag.Name == "-d" || flag.Name == "--days" {
					if flag.Value == nil {
						return nil, fmt.Errorf("flag value required")
					}
					parsedDays, err := strconv.ParseInt(flag.Value.Value, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("flag value must be a valid int")
					}
					days = parsedDays
				}
			}
			return executor.NewStringModel(fmt.Sprintf("the weather for the next %d days is nice", days)), nil
		case "secret":
			return executor.NewStringModel("the secret is " + m.secret), nil
		}
	case "set":
		switch arg {
		case "secret":
			if len(args) < 2 {
				return nil, fmt.Errorf("secret value required")
			}
			secretVal := args[1].Value
			m.secret = secretVal
			return executor.NewStringModel("Secret updated"), nil
		}
	}
	return nil, fmt.Errorf("Invalid input")
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
							CmdMetadata: commandinput.CmdMetadata{
								Level:               1,
								ShowFlagPlaceholder: true,
							},
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
							CmdMetadata: commandinput.CmdMetadata{
								Level:          1,
								PositionalArgs: textInput.NewPositionalArgs("<secret>"),
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
		secret:      "shhh",
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
