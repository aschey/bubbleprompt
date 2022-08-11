package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"runtime"
	"time"

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
	//suggestions []input.Suggestion[cmdMetadata]
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
		txValue := reflect.ValueOf(tx)

		for i := 0; i < txType.NumMethod(); i++ {
			method := txType.Method(i)
			funcN := runtime.FuncForPC(method.Func.Pointer())
			fileName, _ := funcN.FileLine(method.Func.Pointer())
			fset := token.NewFileSet()
			parsedAst, _ := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
			desc := ""
			for _, dec := range parsedAst.Decls {
				if funcDecl, ok := dec.(*ast.FuncDecl); ok {
					if funcDecl.Name.Name == method.Name {
						if funcDecl.Doc != nil {
							desc = funcDecl.Doc.Text()
							if len(desc) > 20 {
								desc = desc[:20]
							}
						}
						break
					}
				}
			}
			// pkg := &ast.Package{
			// 	Name:  "Any",
			// 	Files: make(map[string]*ast.File),
			// }
			// pkg.Files[fileName] = parsedAst

			// importPath, _ := filepath.Abs("/")
			// myDoc := doc.New(pkg, importPath, doc.AllDecls)
			// doc := ""
			// for _, theFunc := range myDoc.Funcs {
			// 	if theFunc.Name == method.Name {
			// 		doc = theFunc.Doc
			// 	}
			// }
			methodVal := txValue.MethodByName(method.Name).Type()
			args := []commandinput.PositionalArg{}
			for j := 0; j < methodVal.NumIn(); j++ {
				param := methodVal.In(j)
				args = append(args, commandinput.NewPositionalArg(param.Name()))
			}

			suggestions = append(suggestions, input.Suggestion[cmdMetadata]{Text: method.Name, Description: desc, Metadata: commandinput.NewCmdMetadata(args, commandinput.Placeholder{})})
		}
		return nil
	})
	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), suggestions), nil
}

func executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "result is " + input, nil
	}), nil
}

func main() {
	config := &flashdb.Config{}
	db, _ := flashdb.New(config)

	var textInput input.Input[cmdMetadata] = commandinput.New[cmdMetadata]()
	completerModel := completerModel{db: db, textInput: textInput.(*commandinput.Model[cmdMetadata])}

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
