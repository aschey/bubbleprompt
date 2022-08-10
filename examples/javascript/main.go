package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/styles"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parser"
	"github.com/aschey/bubbleprompt/input/parserinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dop251/goja"
)

const arrayType = "[]interface {}"
const objectType = "map[string]interface {}"

type model struct {
	prompt prompt.Model[any]
	vm     *vm
}

type completerModel struct {
	textInput   *parserinput.ParserModel[statement]
	suggestions []input.Suggestion[any]
	vm          *vm
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

func (m completerModel) globalSuggestions() []input.Suggestion[any] {
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	vars := m.vm.GlobalObject().Keys()
	suggestions := []input.Suggestion[any]{}
	for _, v := range vars {
		suggestions = append(suggestions, input.Suggestion[any]{Text: v})
	}

	return completers.FilterHasPrefix(currentBeforeCursor, suggestions)
}

func (m completerModel) valueSuggestions(value goja.Value) []input.Suggestion[any] {
	suggestions := []input.Suggestion[any]{}
	if value == nil {
		return suggestions
	}
	strVal := value.String()
	if strVal == "null" || strVal == "undefined" {
		return suggestions
	}
	objectVar := m.vm.ToObject(value)

	_, currentToken := m.textInput.CurrentToken()
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	if currentBeforeCursor == "]" {
		return suggestions
	}

	keyWrap := ""
	datatype := objectVar.ExportType().String()
	// Can't use dot notation with arrays
	if datatype == arrayType && currentBeforeCursor == "." {
		return suggestions
	}

	completable := m.textInput.CompletableTokenBeforeCursor()
	prev := m.textInput.FindLast(func(token parser.Token, symbol string) bool {
		return token.Start < currentToken.Start && symbol != "Whitespace"
	})
	prevToken := ""
	if prev != nil {
		prevToken = prev.Value
	}

	if datatype == objectType && currentBeforeCursor != "." && prevToken != "." && !objectVar.Equals(m.vm.GlobalObject()) {
		keyWrap = `"`
		completable = strings.Trim(completable, `"`)
	}

	for _, key := range objectVar.Keys() {
		suggestions = append(suggestions, input.Suggestion[any]{
			Text:           keyWrap + key + keyWrap,
			CompletionText: key,
		})
	}

	return completers.FilterCompletionTextHasPrefix(completable, suggestions)
}

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[any]) []input.Suggestion[any] {
	parsed := m.textInput.ParsedBeforeCursor()
	if parsed != nil {
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
				json, err := m.vm.ToObject(res).MarshalJSON()
				return string(json), err
			}
		}

		return res.ToString().String(), nil

	}), nil
}

func main() {
	var textInput input.Input[any] = parserinput.NewParserModel[statement](
		parser.NewParticipleParser(participleParser),
		parserinput.WithDelimiterTokens("Punct", "Whitespace", "And", "Or", "Eq"),
		parserinput.WithStyle(styleLexer, *styles.SwapOff),
	)

	vm := newVm()
	_, _ = vm.RunString(`obj = {a: 2, secondVal: 3, blah: {arg: 1, b: '2'}}`)
	_, _ = vm.RunString(`arr = [1, 2, obj]`)

	completerModel := completerModel{
		suggestions: []input.Suggestion[any]{},
		textInput:   textInput.(*parserinput.ParserModel[statement]),
		vm:          vm,
	}

	m := model{prompt: prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
		prompt.WithViewportRenderer[any](),
	), vm: vm}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
