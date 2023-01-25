package cmdtestapp

import (
	"fmt"
	"os"
	"testing"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/formatter"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
)

type cmdMetadata = commandinput.CommandMetadata

type model struct {
	suggestions []suggestion.Suggestion[cmdMetadata]
	textInput   *commandinput.Model[cmdMetadata]
	inc         int
}

type changeTextMsg struct{}

func suggestions(textInput *commandinput.Model[cmdMetadata]) []suggestion.Suggestion[cmdMetadata] {
	return []suggestion.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test desc", Metadata: commandinput.CommandMetadata{
			PositionalArgs: textInput.NewPositionalArgs(
				"[test placeholder1]",
				"[test placeholder2]",
			),
		}},
		{Text: "second-option", Description: "test desc2", Metadata: commandinput.CommandMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[test placeholder]"),
		}},
		{Text: "third-option", Description: "test desc3", Metadata: commandinput.CommandMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[flags]"),
		}},
		{Text: "fourth-option", Description: "test desc4"},
		{Text: "fifth-option", Description: "test desc5"},
		{Text: "sixth-option", Description: "test desc6"},
		{Text: "seventh-option", SuggestionText: "suggestion text", Description: "test desc7"}}
}

func secondLevelSuggestions(
	textInput *commandinput.Model[cmdMetadata],
) []suggestion.Suggestion[cmdMetadata] {
	return []suggestion.Suggestion[cmdMetadata]{
		{Text: "second-level", Description: "test desc", Metadata: commandinput.CommandMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[placeholder2]"),
			Level:          1,
		}},
	}
}

var flags = []commandinput.FlagInput{
	{Short: "t", Long: "test", Description: "test flag"},
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[cmdMetadata], tea.Cmd) {
	switch msg.(type) {
	case changeTextMsg:
		m.suggestions[0].Text = "changed text"
	case suggestion.PeriodicCompleterMsg:
		m.suggestions[0].Text = "changed text" + fmt.Sprint(m.inc)
		m.inc++
	}

	return m, nil
}

func (m model) Complete(
	promptModel prompt.Model[cmdMetadata],
) ([]suggestion.Suggestion[cmdMetadata], error) {
	time.Sleep(100 * time.Millisecond)
	suggestions := m.suggestions
	if m.textInput.CommandCompleted() {
		if m.textInput.ParsedValue().Command.Value == suggestions[2].Text {
			return m.textInput.FlagSuggestions(
				m.textInput.CurrentTokenBeforeCursor().Value,
				flags,
				nil,
			), nil
		}
		suggestions = secondLevelSuggestions(m.textInput)
	}
	return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor().Value, suggestions), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[cmdMetadata]) (tea.Model, error) {
	switch input {
	case "error":
		return nil, fmt.Errorf("bad things")
	case "oneshot":
		return executor.NewCmdModel("", tea.Sequence(
			tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return changeTextMsg{} }),
			suggestion.OneShotCompleter(100*time.Millisecond),
		)), nil
	case "periodic":
		return executor.NewCmdModel("", tea.Sequence(
			tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return changeTextMsg{} }),
			suggestion.PeriodicCompleter(100*time.Millisecond),
		)), nil

	}
	selectedSuggestion := promptModel.SuggestionManager().SelectedSuggestion()

	return executor.NewAsyncStringModel(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		if selectedSuggestion == nil {
			return "result is " + input, nil
		}
		return "selected suggestion is " + selectedSuggestion.Text, nil
	}), nil
}

func TestApp(t *testing.T) {
	formatter.DefaultNameForeground = "15"
	formatter.DefaultSelectedNameForeground = "8"

	formatter.DefaultDescriptionForeground = "15"
	formatter.DefaultDescriptionBackground = "13"
	formatter.DefaultSelectedDescriptionForeground = "8"
	formatter.DefaultSelectedDescriptionBackground = "13"

	formatter.DefaultScrollbarColor = "8"
	formatter.DefaultScrollbarThumbColor = "15"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	textInput := commandinput.New[cmdMetadata]()
	m := model{suggestions: suggestions(textInput), textInput: textInput}

	promptModel := prompt.New[cmdMetadata](
		m,
		textInput,
	)

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
