package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/arriqaaq/flashdb"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
)

type cmdMetadata = commandinput.CmdMetadata

type model struct {
	prompt prompt.Model[cmdMetadata]
}

type completerModel struct {
	db        *flashdb.FlashDB
	textInput *commandinput.Model[cmdMetadata]
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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[cmdMetadata]) ([]input.Suggestion[cmdMetadata], error) {
	suggestions := []input.Suggestion[cmdMetadata]{}
	m.db.View(func(tx *flashdb.Tx) error {

		txType := reflect.TypeOf(tx)

		for i := 0; i < txType.NumMethod(); i++ {
			method := txType.Method(i)
			funcN := runtime.FuncForPC(method.Func.Pointer())
			fileName, _ := funcN.FileLine(method.Func.Pointer())
			fset := token.NewFileSet()
			parsedAst, _ := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
			desc := ""
			args := []commandinput.PositionalArg{}
			for _, dec := range parsedAst.Decls {
				if funcDecl, ok := dec.(*ast.FuncDecl); ok {
					if funcDecl.Name.Name == method.Name {
						for _, arg := range funcDecl.Type.Params.List {
							for _, name := range arg.Names {
								args = append(args, commandinput.NewPositionalArg(name.Name))
							}

						}
						if funcDecl.Doc != nil {
							desc = funcDecl.Doc.Text()
							desc = strings.ReplaceAll(desc, "\n", " ")
							if len(desc) > 80 {
								desc = desc[:80]
							}
						}
						break
					}
				}
			}

			suggestions = append(suggestions, input.Suggestion[cmdMetadata]{Text: method.Name, Description: desc, Metadata: commandinput.NewCmdMetadata(args, commandinput.Placeholder{})})
		}
		return nil
	})
	if m.textInput.CommandCompleted() {
		return nil, nil
	}
	return completers.FilterHasPrefix(m.textInput.CommandBeforeCursor(), suggestions), nil
}

func (m completerModel) executor(input string) (tea.Model, error) {
	if len(input) > 0 {
		input = strings.ToUpper(string(input[0])) + input[1:]
	}
	outStr := ""
	err := m.db.Update(func(tx *flashdb.Tx) error {
		params := strings.Split(input, " ")
		method, found := reflect.TypeOf(tx).MethodByName(params[0])
		if !found {
			return fmt.Errorf("command not found")
		}

		expectedParams := method.Type.NumIn()
		if len(params) != expectedParams {
			// Subtract one for the tx object
			return fmt.Errorf("expected %d params but got %d", expectedParams-1, len(params)-1)
		}

		paramVals := []reflect.Value{reflect.ValueOf(tx)}
		if len(params) > 1 {
			for i, p := range params[1:] {
				var reflectVal any
				var err error
				methodParam := method.Type.In(i + 1)
				switch methodParam.Kind() {
				case reflect.Int64:
					reflectVal, err = strconv.ParseInt(p, 10, 64)
				case reflect.Float64:
					reflectVal, err = strconv.ParseFloat(p, 64)
				case reflect.String:
					reflectVal = p
				}
				if err != nil {
					return err
				}
				paramVals = append(paramVals, reflect.ValueOf(reflectVal))

			}
		}

		out := method.Func.Call(paramVals)
		retVals := []string{}
		for _, outVal := range out {
			if outVal.CanInterface() {
				iface := outVal.Interface()
				if iface == nil {
					outStr = ""
				}
				switch ifaceVal := iface.(type) {
				case error:
					retVals = append(retVals, ifaceVal.Error())
				case []string:
					retVals = append(retVals, strings.Join(ifaceVal, ","))
				case string:
					retVals = append(retVals, ifaceVal)
				case bool:
					retVals = append(retVals, strconv.FormatBool(ifaceVal))
				case int64:
					retVals = append(retVals, strconv.FormatInt(ifaceVal, 10))
				}
			} else {
				retVals = append(retVals, outVal.String())
			}

		}
		outStr = strings.Join(retVals, " ")
		return nil
	})

	return executors.NewStringModel(outStr), err
}

func main() {
	config := &flashdb.Config{}
	db, _ := flashdb.New(config)

	var textInput input.Input[cmdMetadata] = commandinput.New[cmdMetadata]()
	completerModel := completerModel{db: db, textInput: textInput.(*commandinput.Model[cmdMetadata])}

	m := model{prompt: prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
