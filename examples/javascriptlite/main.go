package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parserinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dop251/goja"
)

type model struct {
	prompt prompt.Model[statement]
	vm     *goja.Runtime
}

type completerModel struct {
	textInput   *parserinput.Model[statement]
	suggestions []input.Suggestion[statement]
	vm          *goja.Runtime
}

// func (p statement) CurrentToken() string {
// 	// if len(p.Parts) == 0 {
// 	// 	return ""
// 	// }
// 	return "" //p.Parts[len(p.Parts)-1].Obj
// }

var lex = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "whitespace", Pattern: `\s+`},
	{Name: "String", Pattern: `"([^"]*")|('[^']*')`},
	{Name: "And", Pattern: `&&`},
	{Name: "Or", Pattern: `\|\|`},
	{Name: "Eq", Pattern: `===?`},
	{Name: "Number", Pattern: `[0-9]+(\.[0-9]*)*`},
	{Name: "Punct", Pattern: `[-\[!@#$%^&*()+_=\{\}\|:;"'<,>.?/\]|]`},
	{Name: "Ident", Pattern: `[_a-zA-Z]+[_a-zA-Z0-9]*`},
})

var parser = participle.MustBuild[statement](participle.Lexer(lex), participle.UseLookahead(20))

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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[statement]) []input.Suggestion[statement] {
	vars := m.vm.GlobalObject().Keys()
	suggestions := []input.Suggestion[statement]{}
	for _, v := range vars {
		suggestions = append(suggestions, input.Suggestion[statement]{Text: v})
	}
	//println("current", m.textInput.CurrentTokenBeforeCursor())
	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), suggestions)
}

func (m completerModel) executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		res, err := m.vm.RunString(input)
		return res.ToString().String(), err

	}), nil
}

func main() {

	var textInput input.Input[statement] = parserinput.New(parser)
	vm := goja.New()
	_, _ = vm.RunString(`obj = {a: 2, b: 3}`)
	_, _ = vm.RunString(`arr = [1, 2, 3]`)
	vm.GlobalObject().Keys()
	completerModel := completerModel{
		suggestions: []input.Suggestion[statement]{},
		textInput:   textInput.(*parserinput.Model[statement]),
		vm:          vm,
	}

	m := model{prompt: prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	), vm: vm}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
