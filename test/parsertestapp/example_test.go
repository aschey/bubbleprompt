package parsertestapp

import (
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/participle/v2"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/input/parser"
	"github.com/aschey/bubbleprompt/input/parserinput"
	tea "github.com/charmbracelet/bubbletea"
)

var participleParser = participle.MustBuild[Statement](
	participle.UseLookahead(20),
)

type model struct {
	promptModel prompt.Model[any]
}

type Statement struct {
	Words []string `parser:" (@Ident ( ',' @Ident )*)* "`
}

type completerModel struct {
	textInput   *parserinput.ParserModel[Statement]
	suggestions []input.Suggestion[any]
}

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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
	current := m.textInput.CompletableTokenBeforeCursor()
	suggestions := []input.Suggestion[any]{
		{Text: "abcd"},
		{Text: "def"},
		{Text: "abcdef"},
	}
	return completers.FilterHasPrefix(current, suggestions), nil
}

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[any]) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		err := m.textInput.Error()
		if err != nil {
			return "", err
		}

		return "", nil

	}), nil
}

func TestApp(t *testing.T) {
	input.DefaultNameForeground = "15"
	input.DefaultSelectedNameForeground = "8"

	input.DefaultDescriptionForeground = "15"
	input.DefaultDescriptionBackground = "13"
	input.DefaultSelectedDescriptionForeground = "8"
	input.DefaultSelectedDescriptionBackground = "13"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	prompt.DefaultScrollbarColor = "8"
	prompt.DefaultScrollbarThumbColor = "15"

	textInput := parserinput.NewParserModel[Statement](
		parser.NewParticipleParser(participleParser),
		parserinput.WithDelimiters(","))

	completerModel := completerModel{
		suggestions: []input.Suggestion[any]{},
		textInput:   textInput,
	}

	promptModel, _ := prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
		prompt.WithViewportRenderer[any](),
	)
	m := model{promptModel}

	if _, err := tea.NewProgram(m, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
