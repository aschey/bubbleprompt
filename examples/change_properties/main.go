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
	"github.com/aschey/bubbleprompt/renderer"
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

type appModel struct {
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
}

func (m appModel) Complete(promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	return completer.GetRecursiveCompletions(m.textInput.Tokens(), m.textInput.CursorIndex(), m.suggestions), nil
}

func (m appModel) Execute(input string, promptModel *prompt.Model[cmdMetadata]) (tea.Model, error) {
	parsed := m.textInput.ParsedValue()
	args := parsed.Args
	if len(args) == 0 {
		return nil, fmt.Errorf("At least one argument is required")
	}
	inputFormatters := m.textInput.Formatters()
	promptFormatters := promptModel.Formatters()

	switch parsed.Command.Value() {
	case "cursor-mode":
		switch args[0].Value() {
		case "blink":
			return executor.NewCmdModel("blinking cursor", m.textInput.SetCursorMode(textinput.CursorBlink)), nil
		case "static":
			return executor.NewCmdModel("static cursor", m.textInput.SetCursorMode(textinput.CursorStatic)), nil
		case "hide":
			return executor.NewCmdModel("blinking cursor", m.textInput.SetCursorMode(textinput.CursorHide)), nil
		}
	case "suggestion":
		if len(args) < 2 {
			return nil, fmt.Errorf("At least two arguments are required")
		}
		color := args[1].Value()

		switch args[0].Value() {
		case "name":
			promptFormatters.Name.Style = promptFormatters.Name.Style.Background(lipgloss.Color(color))
		case "description":
			promptFormatters.Description.Style = promptFormatters.Description.Style.Background(lipgloss.Color(color))
		}

	case "input":
		if len(args) < 2 {
			return nil, fmt.Errorf("At least two arguments are required")
		}
		color := args[1].Value()

		switch args[0].Value() {
		case "selected":
			inputFormatters.SelectedText = inputFormatters.SelectedText.Foreground(lipgloss.Color(color))
		case "cursor":
			inputFormatters.Cursor = inputFormatters.Cursor.Foreground(lipgloss.Color(color))
		}

	case "prompt":
		promptValue := args[0].Value()
		m.textInput.SetPrompt(promptValue + " ")
		if len(args) > 1 {
			inputFormatters.Prompt = inputFormatters.Prompt.Foreground(lipgloss.Color(args[1].Value()))
		}

	case "max-suggestions":
		maxSuggestions, err := strconv.ParseInt(args[0].Value(), 10, 64)
		if err != nil {
			return nil, err
		}
		promptModel.SetMaxSuggestions(int(maxSuggestions))

	case "renderer":
		switch args[0].Value() {
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

func (m appModel) Update(msg tea.Msg) (prompt.AppModel[cmdMetadata], tea.Cmd) {
	return m, nil
}

func main() {
	textInput := commandinput.New[cmdMetadata]()

	commandMetadata := commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<command>"))
	colorMetadata := commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<color>"))
	colorMetadata.Level = 1

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
		{
			Text:        "suggestion",
			Description: "set suggestion styles",
			Metadata: cmdMetadata{
				CmdMetadata: commandMetadata,
				children: []input.Suggestion[cmdMetadata]{
					{
						Text:        "name",
						Description: "set suggestion name background",
						Metadata: cmdMetadata{
							CmdMetadata: colorMetadata,
						},
					},
					{
						Text:        "description",
						Description: "set suggestion description background",
						Metadata: cmdMetadata{
							CmdMetadata: colorMetadata,
						},
					},
				},
			},
		},
		{
			Text:        "input",
			Description: "set input style",
			Metadata: cmdMetadata{
				CmdMetadata: commandMetadata,
				children: []input.Suggestion[cmdMetadata]{
					{
						Text:        "selected",
						Description: "set selected suggestion foreground",
						Metadata: cmdMetadata{
							CmdMetadata: colorMetadata,
						},
					},
					{
						Text:        "cursor",
						Description: "set cursor foreground",
						Metadata: cmdMetadata{
							CmdMetadata: colorMetadata,
						},
					},
				},
			},
		},
		{
			Text:        "prompt",
			Description: "set prompt text and foreground",
			Metadata: cmdMetadata{
				CmdMetadata: commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<value>"), textInput.NewPositionalArg("[color]")),
			},
		},
		{
			Text:        "max-suggestions",
			Description: "set max suggestions",
			Metadata: cmdMetadata{
				CmdMetadata: commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<number of suggestions>")),
			},
		},
		{
			Text:        "renderer",
			Description: "change the renderer",
			Metadata: cmdMetadata{
				CmdMetadata: commandMetadata,
				children: []input.Suggestion[cmdMetadata]{
					{Text: "unmanaged", Description: "use the unmanaged renderer"},
					{Text: "viewport", Description: "use the viewport renderer"},
				},
			},
		},
	}
	appModel := appModel{
		suggestions: suggestions,
		textInput:   textInput,
	}

	promptModel, err := prompt.New[cmdMetadata](
		appModel,
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
