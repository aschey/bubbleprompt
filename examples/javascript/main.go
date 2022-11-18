package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/styles"
	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/editor/parser"
	"github.com/aschey/bubbleprompt/editor/parserinput"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/renderer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dop251/goja"
)

const arrayType = "array"
const objectType = "object"
const stringType = "string"

type model struct {
	textInput   *parserinput.ParserModel[any, statement]
	suggestions []editor.Suggestion[any]
	vm          *vm
}

func (m model) globalSuggestions() []editor.Suggestion[any] {
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	vars := m.vm.GlobalObject().Keys()
	suggestions := []editor.Suggestion[any]{}
	for _, v := range vars {
		suggestions = append(suggestions, editor.Suggestion[any]{Text: v})
	}

	return completer.FilterHasPrefix(currentBeforeCursor, suggestions)
}

func (m model) valueSuggestions(value goja.Value) []editor.Suggestion[any] {
	if value == nil {
		return nil
	}
	strVal := value.String()
	if strVal == "null" || strVal == "undefined" {
		return nil
	}
	objectVar := m.vm.ToObject(value)

	currentToken := m.textInput.CurrentToken()
	currentBeforeCursor := m.textInput.CurrentTokenBeforeCursor()
	if currentBeforeCursor == "]" {
		return nil
	}

	keyWrap := ""
	datatype := strings.ToLower(objectVar.ClassName())

	// Can't use dot notation with arrays
	if (datatype == arrayType || datatype == stringType) && currentBeforeCursor == "." {
		return nil
	}

	completable := m.textInput.CompletableTokenBeforeCursor()
	prev := m.textInput.FindLast(func(token editor.Token, symbol string) bool {
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

	suggestions := []editor.Suggestion[any]{}
	for _, key := range objectVar.Keys() {
		suggestions = append(suggestions, editor.Suggestion[any]{
			Text:           keyWrap + key + keyWrap,
			SuggestionText: key,
		})
	}

	return completer.FilterHasPrefix(completable, suggestions)
}

func (m model) Complete(promptModel prompt.Model[any]) ([]editor.Suggestion[any], error) {
	parsed, err := m.textInput.ParsedBeforeCursor()
	if err != nil {
		return nil, err
	}
	if parsed != nil {
		value := m.evaluateStatement(*parsed)
		return m.valueSuggestions(value), nil
	}

	return m.globalSuggestions(), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	return executor.NewAsyncStringModel(func() (string, error) {
		err := m.textInput.Error()
		if err != nil {
			return "", err
		}

		res, err := m.vm.RunString(input)
		if res == nil || err != nil {
			return "", err
		}

		object := m.vm.ToObject(res)
		datatype := strings.ToLower(object.ClassName())

		switch datatype {
		case arrayType, objectType:
			jsonData, err := object.MarshalJSON()
			if err != nil {
				return "", err
			}
			buf := bytes.Buffer{}
			err = json.Indent(&buf, jsonData, "", "  ")
			if err != nil {
				return "", err
			}
			return m.textInput.FormatText(buf.String()), err
		default:
			return m.textInput.FormatText(res.ToString().String()), nil
		}

	}), nil
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	return m, nil
}

func main() {
	textInput := parserinput.NewParserModel[any, statement](
		parser.NewParticipleParser(participleParser),
		parserinput.WithDelimiterTokens[any]("Punct", "Whitespace", "And", "Or", "Eq"),
		parserinput.WithFormatter[any](parser.NewChromaFormatter(styles.SwapOff, styleLexer)),
	)

	vm := newVm()
	_, _ = vm.RunString(`pizza = {mushroom: 'magic', cheese: true, meat: {pepperoni: 1, sausage: 2 }}`)
	_, _ = vm.RunString(`food = ['hummus', 'wine', {pizza: pizza}]`)

	model := model{
		suggestions: []editor.Suggestion[any]{},
		textInput:   textInput,
		vm:          vm,
	}

	promptModel, err := prompt.New[any](
		model,
		textInput,
		prompt.WithViewportRenderer[any](renderer.ViewportOffset{HeightOffset: 1}),
	)
	if err != nil {
		panic(err)
	}

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
