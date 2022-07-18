package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

func (m completerModel) objSuggestions(parent *goja.Object, accessor *propAccessor) []input.Suggestion[tokenMetadata] {
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	varName := accessor.Identifier.Variable
	var obj *goja.Object
	if parent == nil {
		obj = m.vm.Get(varName).ToObject(m.vm)
	} else {
		obj = parent.Get(varName).ToObject(m.vm)
	}
	if accessor.Identifier.Accessor != nil {
		return m.expressionSuggestions(obj, accessor.Identifier.Accessor)
	}
	if obj.ExportType().String() == arrayType {
		return []input.Suggestion[tokenMetadata]{}
	}
	fields := obj.Keys()
	suggestions := []input.Suggestion[tokenMetadata]{}
	skipPrevious := false
	if accessor.Prop == nil {
		currentBeforeCursor = ""
		skipPrevious = true
	}

	for _, f := range fields {
		datatype := obj.Get(f).ExportType().String()
		switch datatype {
		case objectType:
			if accessor.Accessor != nil {
				return m.objSuggestions(obj, accessor.Accessor)
			}
		}

		suggestions = append(suggestions, input.Suggestion[tokenMetadata]{
			Text:     f,
			Metadata: tokenMetadata{skipPrevious},
		})

	}
	return completers.FilterHasPrefix(currentBeforeCursor, suggestions)
}

func (m completerModel) accessorSuggestions(variable *identifier, filterText string, skipPrevious bool) []input.Suggestion[tokenMetadata] {
	curVar := m.vm.Get(variable.Variable)
	if curVar == nil {
		return []input.Suggestion[tokenMetadata]{}
	}
	switch curVar.ExportType().String() {
	case arrayType:
		return m.accessorSuggestionsHelper(curVar, filterText, skipPrevious, nil)

	case objectType:
		return m.accessorSuggestionsHelper(curVar, filterText, skipPrevious, func(key string) string { return `"` + key + `"` })
	}
	return []input.Suggestion[tokenMetadata]{}
}

func (m completerModel) accessorSuggestionsHelper(curVar goja.Value, filterText string, skipPrevious bool, keyFormatter func(key string) string) []input.Suggestion[tokenMetadata] {
	objectVar := curVar.ToObject(m.vm)
	suggestions := []input.Suggestion[tokenMetadata]{}
	for _, key := range objectVar.Keys() {
		if keyFormatter != nil {
			key = keyFormatter(key)
		}
		suggestions = append(suggestions, input.Suggestion[tokenMetadata]{
			Text:     key,
			Metadata: tokenMetadata{skipPrevious},
		})
	}
	return completers.FilterHasPrefix(filterText, suggestions)
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

func (m completerModel) expressionSuggestions(parent *goja.Object, expression *expression) []input.Suggestion[tokenMetadata] {
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	switch {
	case expression.PropAccessor != nil:
		return m.objSuggestions(parent, expression.PropAccessor)
	case expression.Token != nil:
		token := expression.Token
		switch {
		case token.Variable != nil:
			variable := token.Variable
			switch {
			case currentBeforeCursor == "[":
				return m.accessorSuggestions(variable, "", true)
			case variable.Accessor != nil:
				return m.accessorSuggestions(variable, currentBeforeCursor, false)
			case currentBeforeCursor == "":
				return []input.Suggestion[tokenMetadata]{}
			case parent == nil:
				return m.globalSuggestions()
			}
		case token.Literal != nil:
			literal := token.Literal
			literalVal := ""
			switch {
			case literal.Str != nil:
				literalVal = strings.ReplaceAll(*token.Literal.Str, `"`, "")
			case literal.Boolean != nil:
				literalVal = strconv.FormatBool(*literal.Boolean)
			case literal.Number != nil:
				literalVal = strconv.FormatFloat(*literal.Number, 'f', 64, 64)
			}

			return m.accessorSuggestionsHelper(parent.Get(literalVal), "", true, nil)
		}

	}
	return []input.Suggestion[tokenMetadata]{}
}

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[tokenMetadata]) []input.Suggestion[tokenMetadata] {
	parsed := m.textInput.Parsed()
	if parsed != nil {
		switch {
		case parsed.Expression != nil:
			return m.expressionSuggestions(nil, parsed.Expression)
		case parsed.Assignment != nil:
			return []input.Suggestion[tokenMetadata]{}
		}
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
