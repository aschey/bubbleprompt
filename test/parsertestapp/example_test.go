package parsertestapp

import (
	"fmt"
	"os"
	"testing"
	"text/scanner"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parserinput"
	tea "github.com/charmbracelet/bubbletea"
)

var parser = participle.MustBuild[Statement](
	participle.UseLookahead(20),
	participle.Lexer(lexer.NewTextScannerLexer(func(scanner *scanner.Scanner) { scanner.Whitespace = 0 })),
)

type model struct {
	prompt prompt.Model[any]
}

type Statement struct {
	Words []string `parser:" (@Ident ( ',' @Ident )*)* "`
}

type completerModel struct {
	textInput   *parserinput.ParserModel[Statement]
	suggestions []input.Suggestion[any]
}

func (m model) Init() tea.Cmd {
	return m.prompt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.prompt.Update(msg)
	m.prompt = p
	return m, cmd
}

func (m model) View() string {
	return m.prompt.View()
}

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[any]) []input.Suggestion[any] {
	current := m.textInput.CurrentTokenBeforeCursor()
	suggestions := []input.Suggestion[any]{
		{Text: "abc"},
		{Text: "def"},
	}
	return completers.FilterHasPrefix(current, suggestions)
}

func (m completerModel) executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		err := m.textInput.Error()
		if err != nil {
			return "", err
		}

		return "", nil

	}), nil
}

func TestApp(t *testing.T) {
	var textInput input.Input[any] = parserinput.NewParserModel(parser)

	completerModel := completerModel{
		suggestions: []input.Suggestion[any]{},
		textInput:   textInput.(*parserinput.ParserModel[Statement]),
	}

	m := model{prompt: prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
		prompt.WithViewportRenderer[any](),
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
