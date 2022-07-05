package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/participle/v2"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parserinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt prompt.Model[pathExpr]
}

type completerModel struct {
	textInput   *parserinput.Model[pathExpr]
	suggestions []input.Suggestion[pathExpr]
}

type pathExpr struct {
	Parts []part `parser:"@@ ( '.' @@ )*"`
}

type part struct {
	Obj string `parser:"@Ident"`
	Acc []acc  `parser:"('[' @@ ']')*"`
}

type acc struct {
	Name  *string `parser:"@(String|Char|RawString)"`
	Index *int    `parser:"| @Int"`
}

func (p pathExpr) CurrentToken() string {
	if len(p.Parts) == 0 {
		return ""
	}
	return p.Parts[len(p.Parts)-1].Obj
}

var parser = participle.MustBuild[pathExpr]()

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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[pathExpr]) []input.Suggestion[pathExpr] {
	time.Sleep(100 * time.Millisecond)
	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions)
}

func executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "result is " + input
	}), nil
}

func main() {
	suggestions := []input.Suggestion[pathExpr]{
		{Text: "obj1"},
		{Text: "obj2"},
	}
	var textInput input.Input[pathExpr] = parserinput.New(parser)
	completerModel := completerModel{suggestions: suggestions, textInput: textInput.(*parserinput.Model[pathExpr])}
	m := model{prompt: prompt.New(
		completerModel.completer,
		executor,
		textInput,
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
