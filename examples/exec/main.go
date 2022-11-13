package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cmdMetadata = commandinput.CmdMetadata

type model struct {
	promptModel prompt.Model[cmdMetadata]
}

type completerModel struct {
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
}

type cmdModel struct {
	cmd *exec.Cmd
	err error
}

func (m cmdModel) Init() tea.Cmd {
	return tea.ExecProcess(m.cmd, func(err error) tea.Msg {
		return processFinishedMsg{err}
	})
}

func (m cmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(processFinishedMsg); ok {
		m.err = msg.err
		return m, tea.Quit

	}
	return m, nil
}

func (m cmdModel) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().Background(lipgloss.Color("9")).Render("Error: " + m.err.Error())
	}
	return ""
}

type processFinishedMsg struct{ err error }

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
	if !m.textInput.CommandCompleted() {
		return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
	}

	cmd := m.textInput.SelectedCommand()
	if cmd != nil && len(cmd.Metadata.PositionalArgs) > 0 && len(m.textInput.CompletedArgsBeforeCursor()) == 0 {
		pathCompleter := completer.PathCompleter[cmdMetadata]{}
		return pathCompleter.Complete(m.textInput.CurrentTokenBeforeCursor()), nil
	}
	return nil, nil
}

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
	allValues := m.textInput.AllValues()
	cmd := ""
	args := []string{}
	if len(allValues) > 0 {
		cmd = allValues[0]
	}
	if len(allValues) > 1 {
		for _, arg := range allValues[1:] {
			args = append(args, strings.Trim(arg, "\""))
		}
	}

	return cmdModel{cmd: exec.Command(cmd, args...)}, nil
}

func main() {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.
		Color("6")).
		Render("Run an external command without exiting bubbleprompt.\nCurrently this only works well with commands that produce fullscreen output."))
	fmt.Println()

	textInput := commandinput.New[cmdMetadata]()
	filenameArg := commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("[filename]"))
	suggestions := []input.Suggestion[cmdMetadata]{
		{Text: "vim", Metadata: filenameArg},
		{Text: "emacs", Metadata: filenameArg},
		{Text: "nano", Metadata: filenameArg},
		{Text: "top"},
		{Text: "htop"},
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
