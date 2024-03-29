package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cmdMetadata = commandinput.CommandMetadata[any]

type model struct {
	suggestions []suggestion.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[any]
	filterer    completer.Filterer[cmdMetadata]
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

func (m model) Complete(
	promptModel prompt.Model[cmdMetadata],
) ([]suggestion.Suggestion[cmdMetadata], error) {
	if !m.textInput.CommandCompleted() {
		return m.filterer.Filter(
			m.textInput.CurrentTokenBeforeCursor().Value,
			m.suggestions,
		), nil
	}

	parsed := m.textInput.ParsedValue()
	if len(parsed.Args) > 0 && len(m.textInput.CompletedArgsBeforeCursor()) == 0 {
		pathCompleter := completer.PathCompleter[cmdMetadata]{}
		return pathCompleter.Complete(m.textInput.CurrentTokenBeforeCursor().Value), nil
	}
	return nil, nil
}

func (m model) Execute(input string, promptModel *prompt.Model[cmdMetadata]) (tea.Model, error) {
	allValues := m.textInput.Values()
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[cmdMetadata], tea.Cmd) {
	return m, nil
}

func main() {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.
		Color("6")).
		Render("Run an external command without exiting bubbleprompt.\nCurrently this only works well with commands that produce fullscreen output."),
	)
	fmt.Println()

	textInput := commandinput.New[any]()
	filenameArg := commandinput.MetadataFromPositionalArgs[any](textInput.NewPositionalArg("[filename]"))
	suggestions := []suggestion.Suggestion[cmdMetadata]{
		{Text: "vim", Metadata: filenameArg},
		{Text: "emacs", Metadata: filenameArg},
		{Text: "nano", Metadata: filenameArg},
		{Text: "top"},
		{Text: "htop"},
	}
	model := model{
		suggestions: suggestions,
		textInput:   textInput,
		filterer:    completer.NewPrefixFilter[cmdMetadata](),
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
