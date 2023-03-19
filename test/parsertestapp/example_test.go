package parsertestapp

import (
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/participle/v2"
	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/input/lexerinput"
	"github.com/aschey/bubbleprompt/input/parserinput"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/aschey/bubbleprompt/suggestion"
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
	suggestions []suggestion.Suggestion[any]
	filterer    completer.Filterer[any]
}

func (m model) Complete(promptModel prompt.Model[any]) ([]suggestion.Suggestion[any], error) {
	current := m.textInput.CompletableTokenBeforeCursor()
	suggestions := []suggestion.Suggestion[any]{
		{Text: "abcd"},
		{Text: "def"},
		{Text: "abcdef"},
	}
	return m.filterer.Filter(current, suggestions), nil
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	return m, nil
}

func TestApp(t *testing.T) {
	suggestion.DefaultNameForeground = "15"
	suggestion.DefaultSelectedNameForeground = "8"

	suggestion.DefaultDescriptionForeground = "15"
	suggestion.DefaultDescriptionBackground = "13"
	suggestion.DefaultSelectedDescriptionForeground = "8"
	suggestion.DefaultSelectedDescriptionBackground = "13"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	suggestion.DefaultScrollbarColor = "8"
	suggestion.DefaultScrollbarThumbColor = "15"

	textInput := parserinput.NewModel[any, Statement](
		parser.NewParticipleParser(participleParser),
		lexerinput.WithDelimiters[any](","))

	model := model{
		suggestions: []suggestion.Suggestion[any]{},
		textInput:   textInput,
		filterer:    completer.NewPrefixFilter[any](),
	}

	promptModel := prompt.New[any](
		model,
		textInput,
		prompt.WithViewportRenderer[any](renderer.WithHeightOffset(1)),
	)

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
