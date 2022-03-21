package commandinput

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
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
	prompt           string
	delimiterRegex   *regexp.Regexp
	stringRegex      *regexp.Regexp
	Args             []Arg
	PromptStyle      lipgloss.Style
	TextStyle        lipgloss.Style
	CursorStyle      lipgloss.Style
	PlaceholderStyle lipgloss.Style
	parser           *participle.Parser
	parsedText       *Statement
}

func New(opts ...Option) *Model {
	textinput := textinput.New()
	textinput.Focus()
	model := &Model{
		textinput:        textinput,
		Placeholder:      "",
		prompt:           "> ",
		PlaceholderStyle: textinput.PlaceholderStyle,
		parsedText:       &Statement{},
		delimiterRegex:   regexp.MustCompile(`\s+`),
		stringRegex:      regexp.MustCompile(`[^\s]+`),
	}
	for _, opt := range opts {
		if err := opt(model); err != nil {
			panic(err)
		}
	}

	model.buildParser()
	return model
}

func (m *Model) buildParser() {
	lexer := lexer.MustSimple([]lexer.Rule{
		{Name: "QuotedString", Pattern: `"[^"]*"`},
		{Name: `String`, Pattern: m.stringRegex.String()},
		{Name: "whitespace", Pattern: m.delimiterRegex.String()},
	})
	parser := participle.MustBuild(&Statement{}, participle.Lexer(lexer))
	m.parser = parser
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) SetDelimiterRegex(delimiterRegex *regexp.Regexp) {
	m.delimiterRegex = delimiterRegex
	m.buildParser()
}

func (m *Model) SetStringRegex(stringRegex *regexp.Regexp) {
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

func (m *Model) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	expr := &Statement{}
	err := m.parser.ParseString("", m.Value(), expr)
	if err != nil {
		fmt.Println(err)
	}

	m.parsedText = expr

	return cmd
}

func (m *Model) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion) tea.Cmd {
	m.Args = []Arg{}

	if suggestion == nil {
		// Didn't find any matching suggestions, reset
		m.Placeholder = ""
	} else {
		m.Placeholder = suggestion.Text
		for _, arg := range suggestion.PositionalArgs {
			m.Args = append(m.Args, Arg{
				Text:             arg.Placeholder,
				PlaceholderStyle: arg.PlaceholderStyle.Style,
				ArgStyle:         arg.ArgStyle.Style})
		}
	}

	return nil
}

func (m *Model) OnSuggestionChanged(suggestion input.Suggestion) {
	tokenOffset := m.TokenOffset()

	text := m.Value()
	if tokenOffset > -1 {
		m.SetValue(text[:tokenOffset] + suggestion.Text)
	} else {
		m.SetValue(suggestion.Text)
	}
}

func (m *Model) CompletionText(text string) string {
	return m.CurrentTokenBeforeCursor()
}

func (m *Model) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *Model) Value() string {
	return m.textinput.Value()
}

func (m *Model) ParsedValue() Statement {
	return *m.parsedText
}

func (m *Model) CommandBeforeCursor() string {
	parsed := m.ParsedValue()
	if m.Cursor() >= len(parsed.Command.Value) {
		return parsed.Command.Value
	}
	return parsed.Command.Value[:m.Cursor()]
}

func (m *Model) SetTextStyle(style lipgloss.Style) {
	m.TextStyle = style
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

func (m Model) AllTokens() []ident {
	tokens := []ident{m.parsedText.Command}
	tokens = append(tokens, m.ParsedValue().Args.Value...)
	return tokens
}

func (m Model) AllValues() []string {
	tokens := m.AllTokens()
	values := []string{}
	for _, t := range tokens {
		values = append(values, t.Value)
	}
	return values
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

func (m *Model) Prompt() string {
	return m.prompt
}

func (m *Model) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m Model) cursorInToken(tokens []ident, pos int) bool {
	cursor := m.Cursor()
	return cursor >= tokens[pos].Pos.Offset && cursor <= tokens[pos].Pos.Offset+len(tokens[pos].Value)
}

func (m Model) TokenOffset() int {
	cursor := m.Cursor()
	tokens := m.AllTokens()
	if len(tokens) > 0 {
		// Check if cursor is at the end
		last := tokens[len(tokens)-1]
		if cursor > last.Pos.Offset+len(last.Value) {
			return cursor
		}
	}

	for i := len(tokens) - 1; i >= 0; i-- {
		if m.cursorInToken(tokens, i) {
			return tokens[i].Pos.Offset
		}
	}

	return -1
}

func (m Model) CurrentTokenBeforeCursor() string {
	cursor := m.Cursor()
	tokens := m.AllTokens()
	for i := len(tokens) - 1; i >= 0; i-- {
		if m.cursorInToken(tokens, i) {
			end := cursor - tokens[i].Pos.Offset
			if end < len(tokens[i].Value) {
				return tokens[i].Value[:end]
			}
			return tokens[i].Value
		}
	}

	return ""
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
	return m.Cursor() > m.parsedText.Command.Pos.Offset+len(m.parsedText.Command.Value)
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

	return m.PromptStyle.Render(m.prompt) + viewBuilder.getView()
}
