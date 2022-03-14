package commandinput

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Arg struct {
	Text             string
	PlaceholderStyle lipgloss.Style
	ArgStyle         lipgloss.Style
	Formatter        func(arg string) string
}

type Model struct {
	textinput        textinput.Model
	Placeholder      string
	Prompt           string
	delimiterRegex   string
	stringRegex      string
	Args             []Arg
	PromptStyle      lipgloss.Style
	TextStyle        lipgloss.Style
	CursorStyle      lipgloss.Style
	PlaceholderStyle lipgloss.Style
	parser           *participle.Parser
	parsedText       *Statement
}

func New() Model {
	textinput := textinput.New()

	m := Model{
		textinput:        textinput,
		Placeholder:      "",
		Prompt:           "> ",
		PlaceholderStyle: textinput.PlaceholderStyle,
	}
	m.buildParser()
	return m
}

func (m *Model) buildParser() {
	delimiterRegex := m.delimiterRegex
	stringRegex := m.stringRegex
	if delimiterRegex == "" {
		delimiterRegex = `\s+`
	}
	if stringRegex == "" {
		stringRegex = `[^\s]+`
	}
	lexer := lexer.MustSimple([]lexer.Rule{
		{Name: "QuotedString", Pattern: `"[^"]*"`},
		{Name: `String`, Pattern: stringRegex},
		{Name: "whitespace", Pattern: delimiterRegex},
	})
	parser := participle.MustBuild(&Statement{}, participle.Lexer(lexer))
	m.parser = parser
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) SetDelimiterRegex(delimiterRegex string) {
	m.delimiterRegex = delimiterRegex
	m.buildParser()
}

func (m *Model) SetStringRegex(stringRegex string) {
	m.stringRegex = stringRegex
	m.buildParser()
}

type Statement struct {
	Pos     lexer.Position
	Command ident `parser:"@@?"`
	Args    args  `parser:"@@"`
}

type args struct {
	Pos   lexer.Position
	Value []ident `parser:"@@*"`
}

type ident struct {
	Pos   lexer.Position
	Value string `parser:"@QuotedString | @String"`
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	expr := &Statement{}
	err := m.parser.ParseString("", m.Value(), expr)
	if err != nil {
		fmt.Println(err)
	}

	m.parsedText = expr

	return m, cmd
}

func (m *Model) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m Model) Value() string {
	return m.textinput.Value()
}

func (m Model) ParsedValue() Statement {
	return *m.parsedText
}

func (m *Model) SetValue(s string) {
	m.textinput.SetValue(s)
	expr := &Statement{}
	err := m.parser.ParseString("", m.Value(), expr)
	if err != nil {
		fmt.Println(err)
	}

	m.parsedText = expr
}

func (m Model) Cursor() int {
	return m.textinput.Cursor()
}

func (m *Model) SetCursor(pos int) {
	m.textinput.SetCursor(pos)
}

func (m Model) Focused() bool {
	return m.textinput.Focused()
}

func (m Model) LastArg() *ident {
	parsed := *m.parsedText
	if len(parsed.Args.Value) == 0 {
		return nil
	}
	return &parsed.Args.Value[len(parsed.Args.Value)-1]
}

func (m Model) CommandCompleted() bool {
	if m.parsedText == nil {
		return false
	}
	return m.Cursor() > len(m.parsedText.Command.Value)
}

func (m *Model) Blur() {
	m.textinput.Blur()
}

func (m Model) View() string {
	viewBuilder := newViewBuilder(m)
	text := m.Value()
	leadingSpace := text[:m.parsedText.Command.Pos.Offset]
	viewBuilder.render(leadingSpace, lipgloss.NewStyle())

	command := m.parsedText.Command.Value
	viewBuilder.render(command, m.TextStyle)

	if strings.HasPrefix(m.Placeholder, m.Value()) && m.Placeholder != command {
		viewBuilder.render(m.Placeholder[len(command):], m.PlaceholderStyle)
	}

	spaceCount := m.parsedText.Args.Pos.Offset - viewBuilder.viewLen()
	if spaceCount > 0 {
		spaceBeforeArgs := text[viewBuilder.viewLen():m.parsedText.Args.Pos.Offset]
		viewBuilder.render(spaceBeforeArgs, lipgloss.NewStyle())
	}

	for i, arg := range m.parsedText.Args.Value {
		space := text[viewBuilder.viewLen():arg.Pos.Offset]
		viewBuilder.render(space, lipgloss.NewStyle())

		argStyle := lipgloss.NewStyle()
		if i < len(m.Args) {
			argStyle = m.Args[i].ArgStyle
		}
		viewBuilder.render(arg.Value, argStyle)
	}

	startPlaceholder := len(m.parsedText.Args.Value)
	if startPlaceholder < len(m.Args) {
		for _, arg := range m.Args[startPlaceholder:] {
			if viewBuilder.last() != ' ' {
				viewBuilder.render(" ", lipgloss.NewStyle())
			}

			viewBuilder.render(arg.Text, arg.PlaceholderStyle)
		}
	}

	textWithoutSpace := strings.TrimRight(m.Value(), " ")
	extraSpace := m.Value()[len(textWithoutSpace):]

	viewBuilder.render(extraSpace, lipgloss.NewStyle())

	return m.PromptStyle.Render(m.Prompt) + viewBuilder.getView()
}
