package parsertestapp

import (
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/participle/v2"
	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/editor/commandinput"
	"github.com/aschey/bubbleprompt/editor/parser"
	"github.com/aschey/bubbleprompt/editor/parserinput"
	"github.com/aschey/bubbleprompt/executor"
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
	textInput   *parserinput.ParserModel[any, Statement]
	suggestions []editor.Suggestion[any]
}

func (m model) Complete(promptModel prompt.Model[any]) ([]editor.Suggestion[any], error) {
	current := m.textInput.CompletableTokenBeforeCursor()
	suggestions := []editor.Suggestion[any]{
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
	editor.DefaultNameForeground = "15"
	editor.DefaultSelectedNameForeground = "8"

	editor.DefaultDescriptionForeground = "15"
	editor.DefaultDescriptionBackground = "13"
	editor.DefaultSelectedDescriptionForeground = "8"
	editor.DefaultSelectedDescriptionBackground = "13"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	editor.DefaultScrollbarColor = "8"
	editor.DefaultScrollbarThumbColor = "15"

	textInput := parserinput.NewParserModel[any, Statement](
		parser.NewParticipleParser(participleParser),
		parserinput.WithDelimiters[any](","))

	model := model{
		suggestions: []editor.Suggestion[any]{},
		textInput:   textInput,
	}

	promptModel, _ := prompt.New[any](
		model,
		textInput,
		prompt.WithViewportRenderer[any](renderer.ViewportOffset{HeightOffset: 1}),
	)

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
