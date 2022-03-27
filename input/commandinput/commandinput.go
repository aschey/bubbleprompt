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

type arg struct {
	text             string
	placeholderStyle lipgloss.Style
	argStyle         lipgloss.Style
	persist          bool
}

type Model struct {
	textinput         textinput.Model
	Placeholder       string
	prompt            string
	delimiterRegex    *regexp.Regexp
	stringRegex       *regexp.Regexp
	args              []arg
	originalArgs      []arg
	selectedCommand   *input.Suggestion
	PromptStyle       lipgloss.Style
	TextStyle         lipgloss.Style
	SelectedTextStyle lipgloss.Style
	CursorStyle       lipgloss.Style
	PlaceholderStyle  lipgloss.Style
	parser            *participle.Parser
	parsedText        *Statement
}

func New(opts ...Option) *Model {
	textinput := textinput.New()
	textinput.Focus()
	model := &Model{
		textinput:         textinput,
		Placeholder:       "",
		prompt:            "> ",
		PlaceholderStyle:  textinput.PlaceholderStyle,
		SelectedTextStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		parsedText:        &Statement{},
		delimiterRegex:    regexp.MustCompile(`\s+`),
		stringRegex:       regexp.MustCompile(`[^\s]+`),
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
	if m.CommandCompleted() {
		// If no suggestions, leave args alone
		if suggestion == nil {
			return nil
		}

		m.args = []arg{}
		m.args = append(m.args, m.originalArgs...)
		// Subtract 1 to get arg position because index 0 is the command itself
		index := m.CurrentTokenPos().Index - 1
		if index < len(m.args) {
			// Replace current arg with the suggestion
			m.args[index] = arg{
				text:             suggestion.Text,
				placeholderStyle: m.PlaceholderStyle,
				argStyle:         m.originalArgs[index].argStyle,
				persist:          true,
			}
		}

	} else {
		m.args = []arg{}
		// Keep original args so we can reset to this state later
		m.originalArgs = []arg{}
		if suggestion == nil {
			// Didn't find any matching suggestions, reset
			m.Placeholder = ""
		} else {
			m.Placeholder = suggestion.Text
			for _, posArg := range suggestion.PositionalArgs {
				newArg := arg{
					text:             posArg.Placeholder,
					placeholderStyle: posArg.PlaceholderStyle.Style,
					argStyle:         posArg.ArgStyle.Style,
					persist:          false,
				}
				m.args = append(m.args, newArg)
				m.originalArgs = append(m.originalArgs, newArg)
			}
		}
	}

	return nil
}

func (m *Model) OnSuggestionChanged(suggestion input.Suggestion) {
	tokenPos := m.CurrentTokenPos()
	if tokenPos.Index == 0 {
		m.selectedCommand = &suggestion
	}

	text := m.Value()
	if tokenPos.Start > -1 {
		m.SetValue(text[:tokenPos.Start] + suggestion.Text + text[tokenPos.End:])
		// Sometimes SetValue moves the cursor to the end of the line so we need to move it back to the current token
		m.SetCursor(tokenPos.Start)
	} else {
		m.SetValue(suggestion.Text)
	}
	// Recalculate token end position after setting the value to the new suggestion
	newEnd := m.CurrentTokenPos().End
	// Move cursor to the end of the token
	m.SetCursor(newEnd - suggestion.CursorOffset)
}

func (m *Model) OnSuggestionUnselected() {
	if !m.CommandCompleted() {
		m.selectedCommand = nil
	}
}

func (m *Model) CompletionText(text string) string {
	expr := &Statement{}
	_ = m.parser.ParseString("", text, expr)
	tokens := m.allTokens(expr)
	return m.currentToken(tokens)
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

func (m *Model) SetValue(s string) {
	m.textinput.SetValue(s)
	expr := &Statement{}
	err := m.parser.ParseString("", m.Value(), expr)
	if err != nil {
		fmt.Println(err)
	}

	m.parsedText = expr
}

func (m *Model) IsDelimiter(s string) bool {
	return m.delimiterRegex.MatchString(s)
}

func (m Model) AllTokens() []ident {
	return m.allTokens(m.parsedText)
}

func (m Model) allTokens(statement *Statement) []ident {
	tokens := []ident{statement.Command}
	tokens = append(tokens, statement.Args.Value...)
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

type TokenPos struct {
	Start int
	End   int
	Index int
}

func (m Model) CurrentTokenPos() TokenPos {
	cursor := m.Cursor()
	tokens := m.AllTokens()
	if len(tokens) > 0 {
		// Check if cursor is at the end
		last := tokens[len(tokens)-1]
		index := len(tokens) - 1
		value := m.Value()
		if cursor > 0 && m.IsDelimiter(string(value[cursor-1])) {
			// Haven't started a new token yet, but we have added a delimiter
			// so we'll consider the current token finished
			index++
		}
		if cursor > last.Pos.Offset+len(last.Value) {
			return TokenPos{
				Start: cursor,
				End:   cursor,
				Index: index,
			}
		}
	}

	for i := len(tokens) - 1; i >= 0; i-- {
		if m.cursorInToken(tokens, i) {
			return TokenPos{
				Start: tokens[i].Pos.Offset,
				End:   tokens[i].Pos.Offset + len(tokens[i].Value),
				Index: i,
			}
		}
	}

	return TokenPos{
		Start: -1,
		End:   -1,
		Index: -1,
	}
}

func (m Model) CurrentTokenBeforeCursor() string {
	tokens := m.AllTokens()
	return m.currentTokenBeforeCursor(tokens)
}

func (m Model) HasArgs() bool {
	return len(m.parsedText.Args.Value) > 0
}

func (m Model) currentTokenBeforeCursor(tokens []ident) string {
	cursor := m.Cursor()
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

func (m Model) CurrentToken() string {
	return m.currentToken(m.AllTokens())
}

func (m Model) currentToken(tokens []ident) string {
	for i := len(tokens) - 1; i >= 0; i-- {
		if m.cursorInToken(tokens, i) {
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

	// Render command
	command := m.parsedText.Command.Value
	if m.selectedCommand == nil {
		viewBuilder.render(command, m.TextStyle)
	} else {
		viewBuilder.render(command, m.SelectedTextStyle)
	}

	// Render prefix
	if strings.HasPrefix(m.Placeholder, m.Value()) && m.Placeholder != command {
		viewBuilder.render(m.Placeholder[len(command):], m.PlaceholderStyle)
	}

	// Render space before args
	spaceCount := m.parsedText.Args.Pos.Offset - viewBuilder.viewLen()
	if spaceCount > 0 {
		spaceBeforeArgs := text[viewBuilder.viewLen():m.parsedText.Args.Pos.Offset]
		viewBuilder.render(spaceBeforeArgs, lipgloss.NewStyle())
	}

	// Render args
	for i, arg := range m.parsedText.Args.Value {
		space := text[viewBuilder.viewLen():arg.Pos.Offset]
		viewBuilder.render(space, lipgloss.NewStyle())

		argStyle := lipgloss.NewStyle()
		if i < len(m.args) {
			argStyle = m.args[i].argStyle
		}
		viewBuilder.render(arg.Value, argStyle)
	}

	// Render current arg if persist == true
	currentArg := len(m.parsedText.Args.Value) - 1
	if currentArg >= 0 && currentArg < len(m.args) {
		arg := m.args[currentArg]
		if arg.persist && strings.HasPrefix(arg.text, m.parsedText.Args.Value[currentArg].Value) {
			tokenPos := len(m.parsedText.Args.Value[currentArg].Value)
			viewBuilder.render(arg.text[tokenPos:], arg.placeholderStyle)
		}
	}

	// Render arg placeholders
	startPlaceholder := len(m.parsedText.Args.Value)
	if startPlaceholder < len(m.args) {
		for _, arg := range m.args[startPlaceholder:] {
			last := viewBuilder.last()
			if last == nil || *last != ' ' {
				viewBuilder.render(" ", lipgloss.NewStyle())
			}
			viewBuilder.render(arg.text, arg.placeholderStyle)
		}
	}

	// Render trailing delimiters
	// Don't need to do this if there's no args because the trailing space before args gets rendered above
	if m.HasArgs() {
		value := m.Value()
		delimMatches := m.delimiterRegex.FindAllStringIndex(value, -1)
		if len(delimMatches) > 0 {
			lastMatch := delimMatches[len(delimMatches)-1]
			if lastMatch[1] == len(value) {
				// Text ends with delimiter, get the length without trailing delimiters
				textLength := len(value[:lastMatch[0]])
				// Render the trailing delimiters
				extraSpace := m.Value()[textLength:]
				viewBuilder.render(extraSpace, lipgloss.NewStyle())
			}
		}
	}

	return m.PromptStyle.Render(m.prompt) + viewBuilder.getView()
}
