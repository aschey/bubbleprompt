package commandinput

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parser"
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

type PositionalArg struct {
	Placeholder      string
	PlaceholderStyle input.Text
	ArgStyle         input.Text
}

type Placeholder struct {
	Text  string
	Style input.Text
}

type Flag struct {
	Short            string
	Long             string
	Placeholder      string
	Description      string
	RequiresArg      bool
	PlaceholderStyle input.Text
}

const DefaultPlaceholderForeground = "14"

func NewPositionalArg(placeholder string) PositionalArg {
	placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultPlaceholderForeground))
	return PositionalArg{
		Placeholder: placeholder,
		PlaceholderStyle: input.Text{
			Style: placeholderStyle,
		},
	}
}

type Model[T CmdMetadataAccessor] struct {
	textinput          textinput.Model
	commandPlaceholder string
	subcommandWithArgs string
	suggestionLevel    int
	prompt             string
	defaultDelimiter   string
	delimiterRegex     *regexp.Regexp
	stringRegex        *regexp.Regexp
	origArgs           []arg
	args               []arg
	hasFlags           bool
	argIndex           int
	selectedCommand    *input.Suggestion[T]
	currentFlag        *input.Suggestion[T]
	PromptStyle        lipgloss.Style
	TextStyle          lipgloss.Style
	SelectedTextStyle  lipgloss.Style
	CursorStyle        lipgloss.Style
	PlaceholderStyle   lipgloss.Style
	parser             parser.Parser[Statement]
	parsedText         *Statement
}

type TokenPos struct {
	Start int
	End   int
	Index int
}

type RoundingBehavior int

const (
	RoundUp RoundingBehavior = iota
	RoundDown
)

const DefaultSelectedTextColor = "10"

func New[T CmdMetadataAccessor](opts ...Option[T]) *Model[T] {
	textinput := textinput.New()
	textinput.Focus()
	model := &Model[T]{
		textinput:          textinput,
		commandPlaceholder: "",
		subcommandWithArgs: "",
		prompt:             "> ",
		PlaceholderStyle:   textinput.PlaceholderStyle,
		SelectedTextStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultSelectedTextColor)),
		parsedText:         &Statement{},
		delimiterRegex:     regexp.MustCompile(`\s+`),
		stringRegex:        regexp.MustCompile(`[^\-\s][^\s]*`),
		defaultDelimiter:   " ",
	}
	for _, opt := range opts {
		if err := opt(model); err != nil {
			panic(err)
		}
	}

	model.buildParser()
	return model
}

func (m *Model[T]) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model[T]) SetDelimiterRegex(delimiterRegex *regexp.Regexp) {
	m.delimiterRegex = delimiterRegex
	m.buildParser()
}

func (m *Model[T]) SetStringRegex(stringRegex *regexp.Regexp) {
	m.stringRegex = stringRegex
	m.buildParser()
}

func (m *Model[T]) SetDefaultDelimiter(defaultDelimiter string) {
	m.defaultDelimiter = defaultDelimiter
}

func (m *Model[T]) ShouldSelectSuggestion(suggestion input.Suggestion[T]) bool {
	currentTokenPos := m.CurrentTokenPos(RoundUp)
	currentToken := m.CurrentToken(RoundUp)
	// Only select if cursor is at the end of the token or the input will cut off the part after the cursor
	return m.Cursor() == currentTokenPos.End && currentToken == suggestion.Text
}

func (m *Model[T]) ShouldUnselectSuggestion(prevText string, msg tea.KeyMsg) bool {
	pos := m.Cursor()
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete:
		return pos < len(prevText) && !m.isDelimiter(string(prevText[pos]))
	case tea.KeyRunes, tea.KeySpace:
		if msg.String() != "=" {
			return true
		}
		token := ""
		if m.Cursor() == len(m.Value()) {
			tokens := m.AllTokens()
			token = tokens[len(tokens)-1].Value
		} else {
			token = m.CurrentToken(RoundDown)
		}
		// Don't unselect if the current token is a flag and we're adding an = delimiter
		return !strings.HasPrefix(token, "-")
	default:
		return true
	}
}

func (m *Model[T]) ShouldClearSuggestions(prevText string, msg tea.KeyMsg) bool {
	return m.isDelimiter(msg.String())
}

func (m *Model[T]) SelectedCommand() *input.Suggestion[T] {
	return m.selectedCommand
}

func (m *Model[T]) ArgsBeforeCursor() []string {
	args := []string{}
	textBeforeCursor := m.Value()[:m.Cursor()]

	expr, _ := m.parser.Parse(textBeforeCursor)

	for _, arg := range expr.Args.Value {
		args = append(args, arg.Value)
	}
	return args
}

func (m *Model[T]) CompletedArgsBeforeCursor() []string {
	args := []string{}
	textBeforeCursor := m.Value()[:m.Cursor()]

	expr, _ := m.parser.Parse(textBeforeCursor)

	for _, arg := range expr.Args.Value {
		args = append(args, arg.Value)
	}

	if len(expr.Flags.Value) == 0 && len(textBeforeCursor) > 0 && !m.delimiterRegex.MatchString(string(textBeforeCursor[len(textBeforeCursor)-1])) {
		if len(args) > 0 {
			args = args[:len(args)-1]
		}

	}
	return args
}

func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	if _, ok := msg.(tea.KeyMsg); ok {
		expr, err := m.parser.Parse(m.Value())
		if err == nil {
			m.parsedText = expr
		}
	}

	return cmd
}

func (m *Model[T]) FlagSuggestions(inputStr string, flags []Flag, suggestionFunc func(Flag) T) []input.Suggestion[T] {
	suggestions := []input.Suggestion[T]{}
	isLong := strings.HasPrefix(inputStr, "--")
	isMulti := !isLong && strings.HasPrefix(inputStr, "-") && len(inputStr) > 1
	tokenIndex := m.CurrentTokenPos(RoundUp).Index
	allTokens := m.AllTokens()
	prevToken := allTokens[tokenIndex-1].Value

	currentIsFlag := false
	currentToken := ""
	if tokenIndex < len(allTokens) {
		currentToken = allTokens[tokenIndex].Value
		currentIsFlag = strings.HasPrefix(currentToken, "-")
	}

	curFlagText := ""
	if isMulti {
		curFlagText = string(inputStr[len(inputStr)-1])
	}

	for _, flag := range flags {
		// Don't show any flag suggestions if the current flag requires an arg unless the user skipped the arg and is now typing another flag that does not require an arg
		if ((isMulti && flag.Short == curFlagText) || prevToken == "-"+flag.Short || prevToken == "--"+flag.Long) && flag.RequiresArg && (!currentIsFlag || currentToken == "-"+flag.Short || currentToken == "--"+flag.Long) {
			return []input.Suggestion[T]{}
		}

		long := "--" + flag.Long
		short := "-" + flag.Short
		if ((isLong || flag.Short == "") && strings.HasPrefix(long, inputStr)) || strings.HasPrefix(short, inputStr) || (isMulti && !flag.RequiresArg) {
			suggestion := input.Suggestion[T]{
				Description: flag.Description,
			}
			if isLong {
				suggestion.Text = long
			} else if isMulti {
				suggestion.Text = flag.Short
			} else {
				suggestion.Text = short
			}

			if suggestionFunc == nil {
				if metadata, ok := (*new(T)).Create([]PositionalArg{}, Placeholder{Text: flag.Placeholder, Style: flag.PlaceholderStyle}).(T); ok {
					suggestion.Metadata = metadata
				}
			} else {
				suggestion.Metadata = suggestionFunc(flag)
			}

			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions
}

func (m *Model[T]) getPosArgs(metadata T) []arg {
	args := []arg{}
	for i := 0; i < metadata.GetLevel(); i++ {
		args = append(args, arg{
			text: strconv.FormatInt(int64(i), 10),
		})
	}
	for _, posArg := range metadata.GetPositionalArgs() {
		args = append(args, arg{
			text:             posArg.Placeholder,
			placeholderStyle: posArg.PlaceholderStyle.Style,
			argStyle:         posArg.ArgStyle.Style,
			persist:          false,
		})
	}
	return args
}

func (m *Model[T]) OnUpdateFinish(msg tea.Msg, suggestion *input.Suggestion[T], isSelected bool) tea.Cmd {
	if m.CommandCompleted() {
		// If no suggestions, leave args alone
		if suggestion == nil {
			// Don't reset current flag yet so we can still render the placeholder until the arg gets typed
			if m.currentFlag != nil && m.currentFlag.Metadata.GetFlagPlaceholder().Text == "" {
				m.currentFlag = nil
			}
			// Clear any temporary placeholders
			m.args = m.origArgs
			return nil
		}

		if strings.HasPrefix(suggestion.Text, "-") {
			m.currentFlag = suggestion
		} else {
			m.hasFlags = suggestion.Metadata.GetHasFlags()
			m.currentFlag = nil
		}
		index := m.CurrentTokenPos(RoundUp).Index

		if len(suggestion.Metadata.GetPositionalArgs()) > 0 || index <= m.argIndex {
			m.args = []arg{}
			m.origArgs = []arg{}
			m.argIndex = index
			m.subcommandWithArgs = suggestion.Text
			m.suggestionLevel = suggestion.Metadata.GetLevel()

			newArgs := m.getPosArgs(suggestion.Metadata)
			m.args = append(m.args, newArgs...)
			m.origArgs = append(m.origArgs, newArgs...)

		} else {
			m.args = append([]arg{}, m.origArgs...)
		}
		argIndex := index - 1
		if argIndex >= 0 && argIndex < len(m.args) && !suggestion.Metadata.GetPreservePlaceholder() {
			// Replace current arg with the suggestion
			m.args[argIndex] = arg{
				text:             suggestion.Text,
				placeholderStyle: m.PlaceholderStyle,
				argStyle:         m.args[argIndex].argStyle,
				persist:          true,
			}
		}
	} else {
		m.args = []arg{}
		m.origArgs = []arg{}
		m.suggestionLevel = 0
		if suggestion == nil {
			// Didn't find any matching suggestions, reset
			m.commandPlaceholder = ""
			m.subcommandWithArgs = ""
		} else {
			if !strings.HasPrefix(suggestion.Text, "-") {
				m.hasFlags = suggestion.Metadata.GetHasFlags()
			}

			m.commandPlaceholder = suggestion.Text
			m.subcommandWithArgs = suggestion.Text

			for _, posArg := range suggestion.Metadata.GetPositionalArgs() {
				newArg := arg{
					text:             posArg.Placeholder,
					placeholderStyle: posArg.PlaceholderStyle.Style,
					argStyle:         posArg.ArgStyle.Style,
					persist:          false,
				}
				m.args = append(m.args, newArg)
				m.origArgs = append(m.origArgs, newArg)
			}
		}
	}

	return nil
}

func (m *Model[T]) OnSuggestionChanged(suggestion input.Suggestion[T]) {
	token := m.CurrentToken(RoundUp)
	tokenPos := m.CurrentTokenPos(RoundUp)
	if tokenPos.Index == 0 {
		m.selectedCommand = &suggestion
	}

	text := m.Value()
	if tokenPos.Start > -1 {
		cursor := m.Cursor()
		if strings.HasPrefix(token, "-") && !strings.HasPrefix(suggestion.Text, "-") {
			// Adding an additional flag to the flag group, don't replace the entire token
			rest := ""
			if cursor < len(text) {
				// Add trailing text if we're not at the end of the line
				rest = text[cursor+1:]
			}
			m.SetValue(text[:cursor] + suggestion.Text + rest)
		} else if strings.HasPrefix(token, "-") && !strings.HasPrefix(token, "--") && len(token) > 2 && suggestion.Metadata.GetFlagPlaceholder().Text == "" {
			// handle multi flag like -ab
			if cursor == tokenPos.Start {
				// If cursor is on the leading dash, replace the first two characters of the token ([-ab]c)
				m.SetValue(text[:cursor] + suggestion.Text + text[cursor+2:])
			} else {
				// If the cursor is after the dash, trim the dash from the suggestion and replace the single character on the cursor
				m.SetValue(text[:cursor] + suggestion.Text[1:] + text[cursor+1:])
			}
		} else {
			m.SetValue(text[:tokenPos.Start] + suggestion.Text + text[tokenPos.End:])
			// Sometimes SetValue moves the cursor to the end of the line so we need to move it back to the current token
			m.SetCursor(len(text[:tokenPos.Start]+suggestion.Text) - suggestion.CursorOffset)
		}

	} else {
		m.SetValue(suggestion.Text)
	}
}

func (m *Model[T]) OnSuggestionUnselected() {
	if !m.CommandCompleted() {
		m.selectedCommand = nil
	}
}

func (m *Model[T]) CompletionText(text string) string {
	expr, _ := m.parser.Parse(text)
	tokens := m.allTokens(expr)
	token := m.currentToken(tokens, RoundUp)

	return token
}

func (m *Model[T]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *Model[T]) Value() string {
	return m.textinput.Value()
}

func (m *Model[T]) ParsedValue() Statement {
	return *m.parsedText
}

func (m *Model[T]) CommandBeforeCursor() string {
	parsed := m.ParsedValue()
	if m.Cursor() >= len(parsed.Command.Value) {
		return parsed.Command.Value
	}
	return parsed.Command.Value[:m.Cursor()]
}

func (m *Model[T]) SetValue(s string) {
	m.textinput.SetValue(s)
	expr, err := m.parser.Parse(m.Value())
	if err != nil {
		fmt.Println(err)
	}

	m.parsedText = expr
}

func (m *Model[T]) ResetValue() {
	m.textinput.SetValue("")
	m.parsedText = &Statement{}
}

func (m *Model[T]) isDelimiter(s string) bool {
	return m.delimiterRegex.MatchString(s)
}

func (m Model[T]) AllTokens() []ident {
	return m.allTokens(m.parsedText)
}

func (m Model[T]) AllTokensBeforeCursor() []ident {
	textBeforeCursor := m.Value()[:m.Cursor()]

	expr, _ := m.parser.Parse(textBeforeCursor)
	return m.allTokens(expr)
}

func (m Model[T]) AllValuesBeforeCursor() []string {
	tokens := m.AllTokensBeforeCursor()
	values := []string{}
	for _, t := range tokens {
		values = append(values, t.Value)
	}
	return values
}

func (m Model[T]) allTokens(statement *Statement) []ident {
	tokens := []ident{statement.Command}
	tokens = append(tokens, statement.Args.Value...)
	for _, flag := range statement.Flags.Value {
		tokens = append(tokens, ident{Pos: flag.Pos, Value: flag.Name})
		if flag.Value != nil {
			tokens = append(tokens, *flag.Value)
		}
	}
	return tokens
}

func (m Model[T]) AllValues() []string {
	tokens := m.AllTokens()
	values := []string{}
	for _, t := range tokens {
		values = append(values, t.Value)
	}
	return values
}

func (m Model[T]) Cursor() int {
	return m.textinput.Cursor()
}

func (m *Model[T]) SetCursor(pos int) {
	m.textinput.SetCursor(pos)
}

func (m Model[T]) Focused() bool {
	return m.textinput.Focused()
}

func (m *Model[T]) Prompt() string {
	return m.prompt
}

func (m *Model[T]) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m Model[T]) cursorInToken(tokens []ident, pos int, roundingBehavior RoundingBehavior) bool {
	cursor := m.Cursor()
	isInToken := cursor >= tokens[pos].Pos.Offset && cursor <= tokens[pos].Pos.Offset+len(tokens[pos].Value)
	if isInToken {
		return true
	}
	if roundingBehavior == RoundDown {
		if pos == len(tokens)-1 {
			return true
		}
		return cursor < tokens[pos+1].Pos.Offset
	} else {
		if pos == 0 {
			return false
		}
		return cursor > tokens[pos-1].Pos.Offset+len(tokens[pos-1].Value) && cursor < tokens[pos].Pos.Offset
	}

}

func (m Model[T]) CurrentTokenPos(roundingBehavior RoundingBehavior) TokenPos {
	return m.currentTokenPos(m.AllTokens(), roundingBehavior)
}

func (m Model[T]) currentTokenPos(tokens []ident, roundingBehavior RoundingBehavior) TokenPos {
	cursor := m.Cursor()
	if len(tokens) > 0 {
		last := tokens[len(tokens)-1]
		index := len(tokens) - 1
		value := m.Value()
		if roundingBehavior == RoundUp && cursor > 0 && (m.isDelimiter(string(value[cursor-1])) || (strings.HasPrefix(last.Value, "-") && string(value[cursor-1]) == "=")) {
			// Haven't started a new token yet, but we have added a delimiter
			// so we'll consider the current token finished
			index++
		}
		// Check if cursor is at the end
		if cursor > last.Pos.Offset+len(last.Value) {
			return TokenPos{
				Start: cursor,
				End:   cursor,
				Index: index,
			}
		}
	}
	for i := 0; i < len(tokens); i++ {
		if m.cursorInToken(tokens, i, roundingBehavior) {
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

func (m Model[T]) CurrentTokenBeforeCursor(roundingBehavior RoundingBehavior) string {
	start := m.CurrentTokenPos(roundingBehavior).Start
	cursor := m.Cursor()
	if start > cursor {
		return ""
	}
	val := m.Value()[start:cursor]
	return val
}

func (m Model[T]) HasArgs() bool {
	return len(m.parsedText.Args.Value) > 0
}

func (m Model[T]) CurrentToken(roundingBehavior RoundingBehavior) string {
	return m.currentToken(m.AllTokens(), roundingBehavior)
}

func (m Model[T]) currentToken(tokens []ident, roundingBehavior RoundingBehavior) string {
	pos := m.currentTokenPos(tokens, roundingBehavior)
	return m.Value()[pos.Start:pos.End]
}

func (m Model[T]) LastArg() *ident {
	parsed := *m.parsedText
	if len(parsed.Args.Value) == 0 {
		return nil
	}
	return &parsed.Args.Value[len(parsed.Args.Value)-1]
}

func (m Model[T]) CommandCompleted() bool {
	if m.parsedText == nil {
		return false
	}
	return m.Cursor() > m.parsedText.Command.Pos.Offset+len(m.parsedText.Command.Value)
}

func (m *Model[T]) Blur() {
	m.textinput.Blur()
}

func (m *Model[T]) OnExecutorFinished() {}

func (m Model[T]) View(viewMode input.ViewMode) string {
	showCursor := !m.textinput.Blink()
	if viewMode == input.Static {
		showCursor = false
	}
	showPlaceholders := viewMode == input.Interactive
	viewBuilder := input.NewViewBuilder(m.Cursor(), m.CursorStyle, m.defaultDelimiter, showCursor)

	// Render command
	command := m.parsedText.Command.Value
	if m.selectedCommand == nil {
		viewBuilder.Render(command, m.parsedText.Command.Pos.Offset, m.TextStyle)
	} else {
		viewBuilder.Render(command, m.parsedText.Command.Pos.Offset, m.SelectedTextStyle)
	}

	// Render prefix
	if showPlaceholders && strings.HasPrefix(m.commandPlaceholder, m.Value()) && m.commandPlaceholder != command {
		viewBuilder.Render(m.commandPlaceholder[len(command):], m.parsedText.Command.Pos.Offset+len(command), m.PlaceholderStyle)
	}

	// Render args
	for i, arg := range m.parsedText.Args.Value {
		argStyle := lipgloss.NewStyle()
		if i < len(m.args) {
			argStyle = m.args[i].argStyle
		}
		viewBuilder.Render(arg.Value, arg.Pos.Offset, argStyle)
	}

	// Render current arg if persist == true
	currentArg := len(m.parsedText.Args.Value) - 1
	if currentArg >= 0 && currentArg < len(m.args) {
		arg := m.args[currentArg]
		argVal := m.parsedText.Args.Value[currentArg].Value
		// Render the rest of the arg placeholder only if the prefix matches
		if arg.persist && strings.HasPrefix(arg.text, argVal) {
			tokenPos := len(argVal)
			viewBuilder.Render(arg.text[tokenPos:], viewBuilder.ViewLen(), arg.placeholderStyle)
		}
	}

	// Render flags
	currentPos := m.CurrentTokenPos(RoundDown).Start
	currentToken := m.CurrentToken(RoundDown)
	for i, flag := range m.parsedText.Flags.Value {
		viewBuilder.Render(flag.Name, flag.Pos.Offset, lipgloss.NewStyle())
		// Render delimiter only once the full flag has been typed
		if m.currentFlag == nil || len(flag.Name) >= len(m.currentFlag.Text) || flag.Value != nil {
			if flag.Delim != nil {
				viewBuilder.Render(flag.Delim.Value, flag.Delim.Pos.Offset, lipgloss.NewStyle())
			}
		}

		if (flag.Pos.Offset != currentPos) && (currentPos < m.Cursor() || i < len(m.parsedText.Flags.Value)-1 || (len(currentToken) > 0 && !strings.HasPrefix(currentToken, "-"))) {
			if flag.Value != nil {
				viewBuilder.Render(flag.Value.Value, flag.Value.Pos.Offset, lipgloss.NewStyle())
			}

		} else {
			// Render current flag
			if m.currentFlag != nil {
				argVal := ""
				if len(m.parsedText.Flags.Value) > 0 {
					argVal = m.parsedText.Flags.Value[len(m.parsedText.Flags.Value)-1].Name
				}

				// Render the rest of the arg placeholder only if the prefix matches
				if showPlaceholders && strings.HasPrefix(m.currentFlag.Text, argVal) {
					tokenPos := len(argVal)
					viewBuilder.Render(m.currentFlag.Text[tokenPos:], viewBuilder.ViewLen(), m.PlaceholderStyle)
				}
			}

			if m.currentFlag != nil && len(m.currentFlag.Metadata.GetFlagPlaceholder().Text) > 0 && flag.Name[len(flag.Name)-1] != '-' {
				if !m.isDelimiter(string(*viewBuilder.Last())) && *viewBuilder.Last() != '=' {
					viewBuilder.Render(m.defaultDelimiter, viewBuilder.ViewLen(), lipgloss.NewStyle())
				}

				if showPlaceholders && flag.Value == nil {
					viewBuilder.RenderPlaceholder(m.currentFlag.Metadata.GetFlagPlaceholder().Text, viewBuilder.ViewLen(), lipgloss.NewStyle().Foreground(lipgloss.Color("14")))
				}

			}
			if flag.Value != nil {
				viewBuilder.Render(flag.Value.Value, flag.Value.Pos.Offset, lipgloss.NewStyle())
			}
		}
	}

	if showPlaceholders && len(m.parsedText.Flags.Value) == 0 && m.hasFlags {
		m.args = append(m.args, arg{text: "[flags]", placeholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("14"))})
	}

	// Render arg placeholders
	startPlaceholder := len(m.parsedText.Args.Value) + len(m.parsedText.Flags.Value)
	// Don't show arg placeholders if the current arg doesn't match the arg we're about to show placeholders for (the user moved the cursor over to the left)
	all := m.AllValues()
	if m.suggestionLevel > len(all)-1 || strings.HasPrefix(m.subcommandWithArgs, all[m.suggestionLevel]) {
		if showPlaceholders && startPlaceholder < len(m.args) {
			for _, arg := range m.args[startPlaceholder:] {
				last := viewBuilder.Last()
				if last == nil || !m.isDelimiter(string(*last)) {
					viewBuilder.Render(m.defaultDelimiter, viewBuilder.ViewLen(), lipgloss.NewStyle())
				}

				viewBuilder.Render(arg.text, viewBuilder.ViewLen(), arg.placeholderStyle)
			}
		}
	}

	// Render trailing text
	value := m.Value()
	viewLen := viewBuilder.ViewLen()
	if len(value) > viewLen {
		viewBuilder.Render(value[viewBuilder.ViewLen():], viewBuilder.ViewLen(), lipgloss.NewStyle())
	}

	return m.PromptStyle.Render(m.prompt) + viewBuilder.GetView()
}
