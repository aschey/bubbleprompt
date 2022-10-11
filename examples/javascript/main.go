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
	promptModel prompt.Model[any]
	vm          *vm
}

type completerModel struct {
	textInput   *parserinput.ParserModel[statement]
	suggestions []input.Suggestion[any]
	vm          *vm
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
	datatype := objectVar.ExportType().String()
	// Can't use dot notation with arrays
	if datatype == arrayType && currentBeforeCursor == "." {
		return nil
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

	suggestions := []input.Suggestion[any]{}
	for _, key := range objectVar.Keys() {
		suggestions = append(suggestions, input.Suggestion[any]{
			Text:           keyWrap + key + keyWrap,
			CompletionText: key,
		})
	}

	return completers.FilterHasPrefix(completable, suggestions)
}

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
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

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[any]) (tea.Model, error) {
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
				if err != nil {
					return "", err
				}
				return m.textInput.FormatText(string(json)), err
			}
		}

		return m.textInput.FormatText(res.ToString().String()), nil

	}), nil
}

func main() {
	textInput := parserinput.NewParserModel[statement](
		parser.NewParticipleParser(participleParser),
		parserinput.WithDelimiterTokens("Punct", "Whitespace", "And", "Or", "Eq"),
		parserinput.WithFormatter(parser.NewChromaFormatter(styles.SwapOff, styleLexer)),
	)

	vm := newVm()
	_, _ = vm.RunString(`obj = {a: 2, secondVal: 3, blah: {arg: 1, b: '2'}}`)
	_, _ = vm.RunString(`arr = [1, 2, obj]`)

	completerModel := completerModel{
		suggestions: []input.Suggestion[any]{},
		textInput:   textInput,
		vm:          vm,
	}

	promptModel, err := prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
		prompt.WithViewportRenderer[any](),
	)
	if err != nil {
		panic(err)
	}
	m := model{promptModel, vm}

	if _, err := tea.NewProgram(m, tea.WithOnQuit(prompt.OnQuit)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
