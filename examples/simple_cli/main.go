package main

import (
	"fmt"
	"os"
	"strconv"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type secretMsg string

type cmdMetadata = commandinput.CommandMetadata[any]

type model struct {
	suggestions        []suggestion.Suggestion[cmdMetadata]
	textInput          *commandinput.Model[any]
	secret             string
	executorValueStyle lipgloss.Style
	filterer           completer.RecursiveFilterer[cmdMetadata]
}

func (m model) Complete(
	promptModel prompt.Model[cmdMetadata],
) ([]suggestion.Suggestion[cmdMetadata], error) {
	parsed := m.textInput.ParsedValue()
	completed := m.textInput.CompletedArgsBeforeCursor()
	if len(completed) == 1 && parsed.Command.Value == "get" && parsed.Args[0].Value == "weather" {
		flags := []commandinput.FlagInput{
			{
				Short:          "d",
				Long:           "days",
				ArgPlaceholder: m.textInput.NewFlagPlaceholder("<int>"),
				Description:    "Forecast days",
			},
		}
		return m.textInput.FlagSuggestions(
			m.textInput.CurrentTokenBeforeCursor().Value,
			flags,
			nil,
		), nil
	}
	return m.filterer.GetRecursiveSuggestions(
		m.textInput.Tokens(),
		m.textInput.CursorIndex(),
		m.suggestions,
	), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[cmdMetadata]) (tea.Model, error) {
	parsed := m.textInput.ParsedValue()
	args := parsed.Args
	flags := parsed.Flags
	if len(args) == 0 {
		return nil, fmt.Errorf("1 argument required")
	}
	arg := args[0]
	switch parsed.Command.Value {
	case "get":
		switch arg.Value {
		case "weather":
			days := "1"
			if len(flags) > 0 {
				flag := flags[0]
				if flag.Name.Value == "-d" || flag.Name.Value == "--days" {
					if flag.Value == nil {
						return nil, fmt.Errorf("flag value required")
					}
					_, err := strconv.ParseInt(flag.Value.Value, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("flag value must be a valid int")
					}
					days = flag.Value.Value
				}
			}
			days = m.executorValueStyle.Render(days)
			value := m.executorValueStyle.Render("cloudy with a chance of meatballs")
			return executor.NewStringModel(
				fmt.Sprintf("weather for the next %s day(s) is: %s", days, value),
			), nil
		case "secret":
			return executor.NewStringModel(
				"the secret is: " + m.executorValueStyle.Render(m.secret),
			), nil
		}
	case "set":
		switch arg.Value {
		case "secret":
			if len(args) < 2 {
				return nil, fmt.Errorf("secret value required")
			}
			secretVal := args[1]

			return executor.NewCmdModel("Secret updated", func() tea.Msg {
				return secretMsg(secretVal.Unquote())
			}), nil
		}
	}
	return nil, fmt.Errorf("Invalid input")
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[cmdMetadata], tea.Cmd) {
	if msg, ok := msg.(secretMsg); ok {
		m.secret = string(msg)
	}
	return m, nil
}

func main() {
	textInput := commandinput.New[any]()
	secretArgs := textInput.NewPositionalArgs("<secret value>")
	secretArgs[0].ArgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	suggestions := []suggestion.Suggestion[cmdMetadata]{
		{
			Text:        "get",
			Description: "retrieve things",
			Metadata: cmdMetadata{
				PositionalArgs: textInput.NewPositionalArgs("<command"),
				Children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "secret",
						Description: "get the secret",
					},
					{
						Text:        "weather",
						Description: "get the weather",
						Metadata: cmdMetadata{
							ShowFlagPlaceholder: true,
						},
					},
				},
			},
		},
		{
			Text:        "set",
			Description: "update things",
			Metadata: cmdMetadata{
				Children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "secret",
						Description: "update the secret",
						Metadata: cmdMetadata{
							PositionalArgs: secretArgs,
						},
					},
				},
			},
		},
	}
	model := model{
		suggestions:        suggestions,
		textInput:          textInput,
		secret:             "hunter2",
		executorValueStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
		filterer:           completer.NewRecursiveFilterer[cmdMetadata](),
	}

	promptModel := prompt.New[cmdMetadata](
		model,
		textInput,
	)

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
