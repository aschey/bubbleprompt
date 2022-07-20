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

const arrayType = "[]interface {}"
const objectType = "map[string]interface {}"

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

func (m completerModel) globalSuggestions() []input.Suggestion[tokenMetadata] {
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	vars := m.vm.GlobalObject().Keys()
	suggestions := []input.Suggestion[tokenMetadata]{}
	for _, v := range vars {
		suggestions = append(suggestions, input.Suggestion[tokenMetadata]{Text: v})
	}

	return completers.FilterHasPrefix(currentBeforeCursor, suggestions)
}

func (m completerModel) valueSuggestions(value goja.Value) []input.Suggestion[tokenMetadata] {
	if value == nil {
		return m.globalSuggestions()
	}
	objectVar := value.ToObject(m.vm)
	suggestions := []input.Suggestion[tokenMetadata]{}
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	_, prev := m.textInput.PreviousToken()
	prevToken := prev.Value
	skipPrevious := false
	keyWrap := ""

	if objectVar.ExportType().String() == objectType {
		keyWrap = `"`
	}
	if currentBeforeCursor == "." || prevToken == "." {
		keyWrap = ""
	}

	if currentBeforeCursor == "." || currentBeforeCursor == "[" {
		currentBeforeCursor = ""
		skipPrevious = true
	}

	for _, key := range objectVar.Keys() {
		suggestions = append(suggestions, input.Suggestion[tokenMetadata]{
			Text:     keyWrap + key + keyWrap,
			Metadata: tokenMetadata{skipPrevious},
		})
	}

	return completers.FilterHasPrefix(currentBeforeCursor, suggestions)
}

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[tokenMetadata]) []input.Suggestion[tokenMetadata] {
	parsed := m.textInput.ParsedBeforeCursor()
	if parsed != nil {
		//repr.Println(parsed)
		value := m.evaluateStatement(*parsed)
		return m.valueSuggestions(value)
	}

	return m.globalSuggestions()
}

func (m completerModel) executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		err := m.textInput.Error()
		if err != nil {
			return "", err
		}
		res, err := m.vm.RunString(input)
		if res == nil || err != nil {
			return "", err
		}

		exportType := res.ExportType()
		if exportType != nil {
			switch exportType.String() {
			case arrayType, objectType:
				json, err := res.ToObject(m.vm).MarshalJSON()
				return string(json), err
			}
		}

		return res.ToString().String(), nil

	}), nil
}

func main() {
	var textInput input.Input[tokenMetadata] = parserinput.New[tokenMetadata](parser)
	vm := goja.New()
	_, _ = vm.RunString(`obj = {a: 2, secondVal: 3, blah: {arg: 1, b: '2'}}`)
	_, _ = vm.RunString(`arr = [1, 2, obj]`)

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
