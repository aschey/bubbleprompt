package parsertestapp

import (
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/participle/v2"
	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/input/lexerinput"
	"github.com/aschey/bubbleprompt/input/parserinput"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/aschey/bubbleprompt/renderer"
	tea "github.com/charmbracelet/bubbletea"
)

var participleParser = participle.MustBuild[Statement](
	participle.UseLookahead(20),
)

type Statement struct {
	Words []string `parser:" (@Ident ( ',' @Ident )*)* "`
}

type model struct {
	textInput   *parserinput.Model[any, Statement]
	suggestions []input.Suggestion[any]
}

func (m model) Complete(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
	current := m.textInput.CompletableTokenBeforeCursor()
	suggestions := []input.Suggestion[any]{
		{Text: "abcd"},
		{Text: "def"},
		{Text: "abcdef"},
	}
	return completer.FilterHasPrefix(current, suggestions), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	return executor.NewAsyncStringModel(func() (string, error) {
		err := m.textInput.Error()
		if err != nil {
			return "", err
		}

		return "", nil

	}), nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	return m, nil
}

func TestApp(t *testing.T) {
	input.DefaultNameForeground = "15"
	input.DefaultSelectedNameForeground = "8"

	input.DefaultDescriptionForeground = "15"
	input.DefaultDescriptionBackground = "13"
	input.DefaultSelectedDescriptionForeground = "8"
	input.DefaultSelectedDescriptionBackground = "13"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	input.DefaultScrollbarColor = "8"
	input.DefaultScrollbarThumbColor = "15"

	textInput := parserinput.NewModel[any, Statement](
		parser.NewParticipleParser(participleParser),
		lexerinput.WithDelimiters[any](","))

	model := model{
		suggestions: []input.Suggestion[any]{},
		textInput:   textInput,
	}

	promptModel := prompt.New[any](
		model,
		textInput,
		prompt.WithViewportRenderer[any](renderer.ViewportOffset{HeightOffset: 1}),
	)

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
