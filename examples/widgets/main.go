package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/editor/commandinput"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cmdMetadata = commandinput.CmdMetadata

type model struct {
	status statusModel
	prompt prompt.Model[cmdMetadata]
}

func (m model) Init() tea.Cmd {
	return m.prompt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	promptModel, cmd := m.prompt.Update(msg)

	cmds = append(cmds, cmd)
	m.prompt = promptModel.(prompt.Model[cmdMetadata])

	m.status, cmd = m.status.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.prompt.View()+"\n", m.status.View(),
	)
}

type textModel struct {
	editor   textinput.Model
	quitting bool
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyEnter:
			m.quitting = true
			return m, tea.Quit
		}
	}
	editor, cmd := m.editor.Update(msg)
	m.editor = editor
	return m, cmd
}

func (m textModel) View() string {
	if m.quitting {
		return "byeeeeeeee"
	}
	return m.editor.View()
}

type statusModel struct {
	statusText string
	style      lipgloss.Style
	size       tea.WindowSizeMsg
}

type updateStatusMsg string

func (m statusModel) View() string {
	return m.style.Width(m.size.Width).PaddingLeft(1).Render(m.statusText)
}

func (m statusModel) Update(msg tea.Msg) (statusModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
	case updateStatusMsg:
		m.statusText = string(msg)
	}

	return m, nil
}

type inputModel struct {
	suggestions []editor.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
	editText    string
	outputStyle lipgloss.Style
}

func (m inputModel) Complete(promptModel prompt.Model[cmdMetadata]) ([]editor.Suggestion[cmdMetadata], error) {
	if m.textInput.CommandCompleted() {
		return nil, nil
	}
	return completer.FilterHasPrefix(m.textInput.ParsedValue().Command.Value(), m.suggestions), nil
}

func (m inputModel) Execute(input string, promptModel *prompt.Model[cmdMetadata]) (tea.Model, error) {
	parsed := m.textInput.ParsedValue()
	switch parsed.Command.Value() {
	case "set-status":
		if len(parsed.Args) == 0 {
			return nil, fmt.Errorf("One arg required")
		}
		arg := parsed.Args[0].Value()
		return executor.NewCmdModel("status updated", func() tea.Msg {
			return updateStatusMsg(arg)
		}), nil
	case "think":
		if len(parsed.Args) == 0 {
			return nil, fmt.Errorf("One arg required")
		}
		arg := parsed.Args[0].Value()
		intArg, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return nil, err
		}
		asyncModel := executor.NewAsyncStringModel(func() (string, error) {
			time.Sleep(time.Second * time.Duration(intArg))
			return "I'm bored", nil
		})
		asyncModel.LoadingMsg = "thinking..."
		return asyncModel, nil
	case "edit":
		ti := textinput.New()
		ti.Placeholder = "Enter some stuff"
		ti.Focus()
		ti.SetValue(m.editText)
		return textModel{editor: ti}, nil
	}
	return nil, nil
}

func (m inputModel) Update(msg tea.Msg) (prompt.InputHandler[cmdMetadata], tea.Cmd) {
	if msg, ok := msg.(prompt.ExecutorFinishedMsg); ok {
		if model, ok := msg.(textModel); ok {
			m.editText = model.editor.Value()
		}
	}
	return m, nil
}

func main() {
	textInput := commandinput.New[cmdMetadata]()
	secondsArg := textInput.NewPositionalArg("<seconds>")
	secondsArg.ArgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	suggestions := []editor.Suggestion[cmdMetadata]{
		{
			Text:        "set-status",
			Description: "set statusbar text",
			Metadata:    commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<text>")),
		},
		{
			Text:        "think",
			Description: "just think for a bit",
			Metadata:    commandinput.MetadataFromPositionalArgs(secondsArg),
		},
		{
			Text:        "edit",
			Description: "edit some text",
		},
	}

	inputModel := inputModel{
		suggestions: suggestions,
		textInput:   textInput,
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}

	promptModel, err := prompt.New[cmdMetadata](
		inputModel,
		textInput,
		prompt.WithViewportRenderer[cmdMetadata](renderer.ViewportOffset{HeightOffset: 2}),
	)
	if err != nil {
		panic(err)
	}

	model := model{
		prompt: promptModel,
		status: statusModel{
			statusText: "all systems go",
			style:      lipgloss.NewStyle().Background(lipgloss.Color("2")).Foreground(lipgloss.Color("15")),
		},
	}

	if _, err := tea.NewProgram(model, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
