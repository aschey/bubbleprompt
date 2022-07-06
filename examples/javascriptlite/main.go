package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parserinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type model struct {
	prompt prompt.Model[statement]
}

type completerModel struct {
	textInput   *parserinput.Model[statement]
	suggestions []input.Suggestion[statement]
}

type statement struct {
	Assignment *assignment `parser:"(@@"`
	Expression *expression `parser:"| @@)?"`
}

type assignment struct {
	Identifier identifier `parser:" @@ '=' "`
	Expression expression `parser:"@@"`
}

type identifier struct {
	Variable string      `parser:"@Ident"`
	Accessor *expression `parser:" ('[' @@ ']')? "`
}

type group struct {
	Expression *expression `parser:"'(' @@ ')'"`
}

type expression struct {
	Array      *array      `parser:"( @@"`
	Object     *object     `parser:"| @@"`
	Group      *group      `parser:"| @@"`
	Prop       *prop       `parser:"| @@"`
	Token      *token      `parser:"| @@)"`
	InfixOp    *infixOp    `parser:"(@@"`
	Expression *expression `parser:"@@)?"`
}

type token struct {
	Literal  *literal    `parser:"@@"`
	Variable *identifier `parser:"| @@"`
}

type keyValuePair struct {
	Key   string     `parser:" @String "`
	Value expression `parser:" ':' @@ "`
}

type object struct {
	Properties *[]keyValuePair `parser:" '{' (@@ (',' @@)*)* '}' "`
}

type prop struct {
	Identifier identifier `parser:" @@ '.' "`
	Prop       string     `parser:"@Ident"`
}

type infixOp struct {
	Op string `parser:" '+' | '-' | '*' | '/' | '||' | '&&' | '==' | '===' "`
}

type array struct {
	Values []expression `parser:" '[' (@@ ( ',' @@ )*)* ']' "`
}

type literal struct {
	Null    *string  `parser:" ( ( 'null' | 'undefined' ) "`
	Boolean *bool    `parser:" | ( 'true' | 'false' ) "`
	Str     *string  `parser:"| @String"`
	Number  *float64 `parser:"| @Number )"`
}

func (p statement) CurrentToken() string {
	// if len(p.Parts) == 0 {
	// 	return ""
	// }
	return "" //p.Parts[len(p.Parts)-1].Obj
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

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[statement]) []input.Suggestion[statement] {
	time.Sleep(100 * time.Millisecond)
	p := m.textInput.Parsed()
	if p != nil {
		spew.Printf("%#v", *p)
	}

	return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions)
}

func executor(input string) (tea.Model, error) {
	return executors.NewAsyncStringModel(func() string {
		time.Sleep(100 * time.Millisecond)
		return "result is " + input
	}), nil
}

func main() {
	suggestions := []input.Suggestion[statement]{
		{Text: "obj1"},
		{Text: "obj2"},
	}
	var textInput input.Input[statement] = parserinput.New(parser)
	completerModel := completerModel{suggestions: suggestions, textInput: textInput.(*parserinput.Model[statement])}
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
