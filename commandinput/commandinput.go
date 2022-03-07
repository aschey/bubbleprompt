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
	Args             []Arg
	PromptStyle      lipgloss.Style
	TextStyle        lipgloss.Style
	CursorStyle      lipgloss.Style
	PlaceholderStyle lipgloss.Style
	parser           *participle.Parser
	parsedText       *statement
}

func New() Model {
	textinput := textinput.New()
	lexer := lexer.MustSimple([]lexer.Rule{
		{Name: "QuotedString", Pattern: `"[^"]*"`},
		{Name: `String`, Pattern: `[^\s]+`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	parser := participle.MustBuild(&statement{}, participle.Lexer(lexer))
	return Model{
		textinput:        textinput,
		Placeholder:      "",
		Prompt:           "> ",
		PlaceholderStyle: textinput.PlaceholderStyle,
		parser:           parser,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type statement struct {
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

	expr := &statement{}
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

func (m *Model) SetValue(s string) {
	m.textinput.SetValue(s)
	expr := &statement{}
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

func (m *Model) Blur() {
	m.textinput.Blur()
}

func (m Model) View() string {
	viewBuilder := newViewBuilder(m)

	leadingSpace := strings.Repeat(" ", m.parsedText.Command.Pos.Offset)
	viewBuilder.render(leadingSpace, lipgloss.NewStyle())

	viewBuilder.render(m.parsedText.Command.Value, m.TextStyle)

	if strings.HasPrefix(m.Placeholder, m.Value()) && m.Placeholder != m.parsedText.Command.Value {
		viewBuilder.render(m.Placeholder[len(m.parsedText.Command.Value):], m.PlaceholderStyle)
	}

	spaceCount := m.parsedText.Args.Pos.Offset - viewBuilder.viewLen
	if spaceCount > 0 {
		spaceBeforeArgs := strings.Repeat(" ", m.parsedText.Args.Pos.Offset-viewBuilder.viewLen)
		viewBuilder.render(spaceBeforeArgs, lipgloss.NewStyle())
	}

	for i, arg := range m.parsedText.Args.Value {
		w := arg.Pos.Offset - viewBuilder.viewLen
		space := strings.Repeat(" ", w)
		viewBuilder.render(space, lipgloss.NewStyle())

		argStyle := lipgloss.NewStyle()
		if i < len(m.Args) {
			argStyle = m.Args[i].ArgStyle
		}
		viewBuilder.render(arg.Value, argStyle)
	}
	textWithoutSpace := strings.TrimRight(m.Value(), " ")
	extraSpace := m.Value()[len(textWithoutSpace):]

	viewBuilder.render(extraSpace, lipgloss.NewStyle())

	return m.PromptStyle.Render(m.Prompt) + viewBuilder.getView()
}
