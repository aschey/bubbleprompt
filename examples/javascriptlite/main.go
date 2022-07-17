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
	prompt prompt.Model[tokenMetadata]
	vm     *goja.Runtime
}

type tokenMetadata struct {
	skipPrevious bool
}

func (t tokenMetadata) SkipPrevious() bool {
	return t.skipPrevious
}

type completerModel struct {
	textInput   *parserinput.Model[tokenMetadata, statement]
	suggestions []input.Suggestion[tokenMetadata]
	vm          *goja.Runtime
}

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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[tokenMetadata]) []input.Suggestion[tokenMetadata] {
	vars := m.vm.GlobalObject().Keys()
	//symbols := lexer.SymbolsByRune(lex)
	suggestions := []input.Suggestion[tokenMetadata]{}
	for _, v := range vars {
		suggestions = append(suggestions, input.Suggestion[tokenMetadata]{Text: v})
	}

	// _, current := m.textInput.CurrentToken()
	// _, prev := m.textInput.PreviousToken()
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()

	parsed := m.textInput.Parsed()
	if parsed != nil {
		if parsed.Expression != nil {
			if parsed.Expression.Prop != nil {
				varName := parsed.Expression.Prop.Identifier.Variable
				fields := m.vm.Get(varName).ToObject(m.vm).Keys()
				suggestions = []input.Suggestion[tokenMetadata]{}
				skipPrevious := false
				if parsed.Expression.Prop.Prop == nil {
					currentBeforeCursor = ""
					skipPrevious = true
				}

				for _, f := range fields {
					suggestions = append(suggestions, input.Suggestion[tokenMetadata]{
						Text:     f,
						Metadata: tokenMetadata{skipPrevious},
					})
				}
			}
		}
	}

	// if current != nil && prev != nil {
	// 	current := *current
	// 	prev := *prev
	// 	if symbols[prev.Type] == "Ident" && symbols[current.Type] == "Punct" {
	// 		currentBeforeCursor = ""
	// 		varName := prev.String()
	// 		if slices.Contains(vars, varName) {
	// 			fields := m.vm.Get(varName).ToObject(m.vm).Keys()

	// 			suggestions = []input.Suggestion[tokenMetadata]{}
	// 			for _, f := range fields {
	// 				suggestions = append(suggestions, input.Suggestion[tokenMetadata]{Text: f})
	// 			}
	// 		}
	// 	} else if  symbols[current.Type] == "Punct" &&
	// }

	return completers.FilterHasPrefix(currentBeforeCursor, suggestions)
}

func (m completerModel) executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		res, err := m.vm.RunString(input)
		return res.ToString().String(), err

	}), nil
}

func main() {

	var textInput input.Input[tokenMetadata] = parserinput.New[tokenMetadata](parser)
	vm := goja.New()
	_, _ = vm.RunString(`obj = {a: 2, secondVal: 3}`)
	_, _ = vm.RunString(`arr = [1, 2, 3]`)
	vm.GlobalObject().Keys()
	completerModel := completerModel{
		suggestions: []input.Suggestion[tokenMetadata]{},
		textInput:   textInput.(*parserinput.Model[tokenMetadata, statement]),
		vm:          vm,
	}

	m := model{prompt: prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
		prompt.WithViewportRenderer[tokenMetadata](),
	), vm: vm}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
