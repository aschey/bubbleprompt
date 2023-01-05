package main

import (
	"fmt"
	"os"
	"strconv"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/formatter"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cmdMetadata struct {
	commandinput.CommandMetadata
	children []suggestion.Suggestion[cmdMetadata]
}

func (c cmdMetadata) Children() []suggestion.Suggestion[cmdMetadata] {
	return c.children
}

type model struct {
	suggestions []suggestion.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
}

func (m model) Complete(promptModel prompt.Model[cmdMetadata]) ([]suggestion.Suggestion[cmdMetadata], error) {
	return completer.GetRecursiveSuggestions(m.textInput.Tokens(), m.textInput.CursorIndex(), m.suggestions), nil
}

func (m model) Execute(inputStr string, promptModel *prompt.Model[cmdMetadata]) (tea.Model, error) {
	parsed := m.textInput.ParsedValue()
	args := parsed.Args
	if len(args) == 0 {
		return nil, fmt.Errorf("At least one argument is required")
	}
	inputFormatters := m.textInput.Formatters()
	promptFormatters := promptModel.Formatters()

	switch parsed.Command.Value {
	case "theme":
		switch args[0].Value {
		case "default":
			promptFormatters = formatter.DefaultFormatters()
			promptModel.SuggestionManager().SetSelectionIndicator("")
		case "minimal":
			promptFormatters = formatter.DefaultFormatters().Minimal()
			promptModel.SuggestionManager().SetSelectionIndicator("> ")
		}
	case "scrollbar":
		switch args[0].Value {
		case "enable":
			promptModel.SuggestionManager().EnableScrollbar()
		case "disable":
			promptModel.SuggestionManager().DisableScrollbar()
		}

	case "cursor-mode":
		switch args[0].Value {
		case "blink":
			return executor.NewCmdModel("blinking cursor", m.textInput.SetCursorMode(cursor.CursorBlink)), nil
		case "static":
			return executor.NewCmdModel("static cursor", m.textInput.SetCursorMode(cursor.CursorStatic)), nil
		case "hide":
			return executor.NewCmdModel("blinking cursor", m.textInput.SetCursorMode(cursor.CursorHide)), nil
		}
	case "suggestion":
		if len(args) < 2 {
			return nil, fmt.Errorf("At least two arguments are required")
		}
		color := args[1].Value

		switch args[0].Value {
		case "name":
			promptFormatters.Name.Style = promptFormatters.Name.Style.Background(lipgloss.Color(color))
		case "description":
			promptFormatters.Description.Style = promptFormatters.Description.Style.Background(lipgloss.Color(color))
		}

	case "input":
		if len(args) < 2 {
			return nil, fmt.Errorf("At least two arguments are required")
		}
		color := args[1].Value

		switch args[0].Value {
		case "selected":
			inputFormatters.SelectedText = inputFormatters.SelectedText.Foreground(lipgloss.Color(color))
		case "cursor":
			inputFormatters.Cursor = inputFormatters.Cursor.Foreground(lipgloss.Color(color))
		}

	case "prompt":
		promptValue := args[0].Value
		m.textInput.SetPrompt(promptValue + " ")
		if len(args) > 1 {
			inputFormatters.Prompt = inputFormatters.Prompt.Foreground(lipgloss.Color(args[1].Value))
		}

	case "max-suggestions":
		maxSuggestions, err := strconv.ParseInt(args[0].Value, 10, 64)
		if err != nil {
			return nil, err
		}
		promptModel.SuggestionManager().SetMaxSuggestions(int(maxSuggestions))

	case "renderer":
		switch args[0].Value {
		case "viewport":
			return executor.NewCmdModel("set viewport renderer", prompt.SetRenderer(renderer.NewViewportRenderer(renderer.ViewportOffset{}), true)), nil
		case "unmanaged":
			return executor.NewCmdModel("set unmanaged renderer", prompt.SetRenderer(renderer.NewUnmanagedRenderer(), true)), nil
		}
	}

	m.textInput.SetFormatters(inputFormatters)
	promptModel.SetFormatters(promptFormatters)
	return executor.NewStringModel("input updated"), nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[cmdMetadata], tea.Cmd) {
	return m, nil
}

func main() {
	textInput := commandinput.New[cmdMetadata]()

	commandMetadata := commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<command>"))
	colorMetadata := cmdMetadata{
		CommandMetadata: commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<color>")),
	}
	colorMetadata.Level = 1

	childMetadata := cmdMetadata{
		CommandMetadata: commandinput.CommandMetadata{
			Level: 1,
		},
	}

	suggestions := []suggestion.Suggestion[cmdMetadata]{
		{
			Text:        "cursor-mode",
			Description: "set the cursor mode",
			Metadata: cmdMetadata{
				CommandMetadata: commandMetadata,
				children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "blink",
						Description: "blinking cursor",
						Metadata:    childMetadata,
					},
					{
						Text:        "static",
						Description: "normal cursor",
						Metadata:    childMetadata,
					},
					{
						Text:        "hide",
						Description: "no cursor",
						Metadata: cmdMetadata{
							CommandMetadata: commandinput.CommandMetadata{
								Level: 1,
							},
						},
					},
				},
			},
		},
		{
			Text:        "suggestion",
			Description: "set suggestion styles",
			Metadata: cmdMetadata{
				CommandMetadata: commandMetadata,
				children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "name",
						Description: "set suggestion name background",
						Metadata:    colorMetadata,
					},
					{
						Text:        "description",
						Description: "set suggestion description background",
						Metadata:    colorMetadata,
					},
				},
			},
		},
		{
			Text:        "input",
			Description: "set input style",
			Metadata: cmdMetadata{
				CommandMetadata: commandMetadata,
				children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "selected",
						Description: "set selected suggestion foreground",
						Metadata:    colorMetadata,
					},
					{
						Text:        "cursor",
						Description: "set cursor foreground",
						Metadata:    colorMetadata,
					},
				},
			},
		},
		{
			Text:        "theme",
			Description: "change theme",
			Metadata: cmdMetadata{
				CommandMetadata: commandMetadata,
				children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "default",
						Description: "enable default theme",
						Metadata:    childMetadata,
					},
					{
						Text:        "minimal",
						Description: "enable the minimal theme",
						Metadata:    childMetadata,
					},
				},
			},
		},
		{
			Text:        "scrollbar",
			Description: "change the scrollbar",
			Metadata: cmdMetadata{
				CommandMetadata: commandMetadata,
				children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "enable",
						Description: "enable the scrollbar",
						Metadata:    childMetadata,
					},
					{
						Text:        "disable",
						Description: "disable the scrollbar",
						Metadata:    childMetadata,
					},
				},
			},
		},
		{
			Text:        "prompt",
			Description: "set prompt text and foreground",
			Metadata: cmdMetadata{
				CommandMetadata: commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<value>"), textInput.NewPositionalArg("[color]")),
			},
		},
		{
			Text:        "max-suggestions",
			Description: "set max suggestions",
			Metadata: cmdMetadata{
				CommandMetadata: commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<number of suggestions>")),
			},
		},
		{
			Text:        "renderer",
			Description: "change the renderer",
			Metadata: cmdMetadata{
				CommandMetadata: commandMetadata,
				children: []suggestion.Suggestion[cmdMetadata]{
					{
						Text:        "unmanaged",
						Description: "use the unmanaged renderer",
						Metadata:    childMetadata,
					},
					{
						Text:        "viewport",
						Description: "use the viewport renderer",
						Metadata:    childMetadata,
					},
				},
			},
		},
	}
	appModel := model{
		suggestions: suggestions,
		textInput:   textInput,
	}

	promptModel := prompt.New[cmdMetadata](
		appModel,
		textInput,
	)

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
