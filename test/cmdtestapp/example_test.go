package cmdtestapp

import (
	"fmt"
	"os"
	"testing"
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
	suggestions []input.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
	doOneshot   bool
	doPeriodic  bool
	inc         int
}

type changeTextMsg struct{}

var suggestions []input.Suggestion[cmdMetadata] = []input.Suggestion[cmdMetadata]{
	{Text: "first-option", Description: "test desc", Metadata: commandinput.CmdMetadata{
		PositionalArgs: []commandinput.PositionalArg{
			commandinput.NewPositionalArg("[test placeholder1]"),
			commandinput.NewPositionalArg("[test placeholder2]"),
		},
	}},
	{Text: "second-option", Description: "test desc2", Metadata: commandinput.CmdMetadata{
		PositionalArgs: []commandinput.PositionalArg{
			commandinput.NewPositionalArg("[test placeholder]"),
		},
	}},
	{Text: "third-option", Description: "test desc3", Metadata: commandinput.CmdMetadata{
		PositionalArgs: []commandinput.PositionalArg{
			commandinput.NewPositionalArg("[flags]"),
		},
	}},
	{Text: "fourth-option", Description: "test desc4"},
	{Text: "fifth-option", Description: "test desc5"},
	{Text: "sixth-option", Description: "test desc6"},
	{Text: "seventh-option", CompletionText: "completion text", Description: "test desc7"}}

var secondLevelSuggestions []input.Suggestion[cmdMetadata] = []input.Suggestion[cmdMetadata]{
	{Text: "second-level", Description: "test desc", Metadata: commandinput.CmdMetadata{
		PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("[placeholder2]")},
		Level:          1,
	}},
}

var flags = []commandinput.Flag{
	{Short: "t", Long: "test", Description: "test flag", RequiresArg: false},
}

func (m *model) Init() tea.Cmd {
	return m.promptModel.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	p, cmd := m.promptModel.Update(msg)
	cmds = append(cmds, cmd)
	m.promptModel = p

	switch msg.(type) {
	case changeTextMsg:
		m.suggestions[0].Text = "changed text"
	case prompt.PeriodicCompleterMsg:
		m.suggestions[0].Text = "changed text" + fmt.Sprint(m.inc)
		m.inc++
	}

	if m.doOneshot {
		m.doOneshot = false
		cmds = append(cmds,
			tea.Sequence(
				tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return changeTextMsg{} }),
				prompt.OneShotCompleter(100*time.Millisecond),
			),
		)
	} else if m.doPeriodic {
		m.doPeriodic = false
		cmds = append(cmds,
			tea.Sequence(
				tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return changeTextMsg{} }),
				prompt.PeriodicCompleter(100*time.Millisecond),
			),
		)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return m.promptModel.View()
}

func (m *model) completer(document prompt.Document, promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	time.Sleep(100 * time.Millisecond)
	suggestions := m.suggestions
	if m.textInput.CommandCompleted() {
		if m.textInput.ParsedValue().Command.Value == suggestions[2].Text {
			return m.textInput.FlagSuggestions(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), flags, nil), nil
		}
		suggestions = secondLevelSuggestions
	}
	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), suggestions), nil
}

func (m *model) executor(input string, selectedSuggestion *input.Suggestion[cmdMetadata]) (tea.Model, error) {
	switch input {
	case "error":
		return nil, fmt.Errorf("bad things")
	case "oneshot":
		m.doOneshot = true
		return executors.NewStringModel(""), nil
	case "periodic":
		m.doPeriodic = true
		return executors.NewStringModel(""), nil
	}

	return executors.NewAsyncStringModel(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		if selectedSuggestion == nil {
			return "result is " + input, nil
		}
		return "selected suggestion is " + selectedSuggestion.Text, nil
	}), nil
}

func TestApp(t *testing.T) {
	input.DefaultNameForeground = "15"
	input.DefaultSelectedNameForeground = "8"

	input.DefaultDescriptionForeground = "15"
	input.DefaultDescriptionBackground = "13"
	input.DefaultSelectedDescriptionForeground = "8"
	input.DefaultSelectedDescriptionBackground = "13"

	prompt.DefaultScrollbarColor = "8"
	prompt.DefaultScrollbarThumbColor = "15"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	textInput := commandinput.New[cmdMetadata]()
	m := model{suggestions: suggestions, textInput: textInput}

	promptModel, _ := prompt.New(
		m.completer,
		m.executor,
		textInput,
	)
	m.promptModel = promptModel

	if _, err := tea.NewProgram(&m, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
