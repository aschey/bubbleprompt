package commandinput

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parser"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type arg struct {
	text             string
	placeholderStyle lipgloss.Style
	argStyle         lipgloss.Style
	persist          bool
}

type PositionalArg struct {
	placeholder      string
	PlaceholderStyle lipgloss.Style
	ArgStyle         lipgloss.Style
}

func (p PositionalArg) Placeholder() string {
	return p.placeholder
}

type FlagPlaceholder struct {
	text  string
	Style lipgloss.Style
}

type Flag struct {
	Short       string
	Long        string
	Placeholder FlagPlaceholder
	Description string
}

func (f Flag) RequiresArg() bool {
	return len(f.Placeholder.text) > 0
}

type Model[T CmdMetadataAccessor] struct {
	textinput           textinput.Model
	commandPlaceholder  []rune
	subcommandWithArgs  string
	suggestionLevel     int
	prompt              string
	defaultDelimiter    string
	delimiterRegex      *regexp.Regexp
	origArgs            []arg
	args                []arg
	showFlagPlaceholder bool
	argIndex            int
	selectedCommand     *input.Suggestion[T]
	currentFlag         *input.Suggestion[T]
	formatters          Formatters
	parser              parser.Parser[Statement]
	parsedText          *Statement
}

type TokenPos struct {
	Start int
	End   int
	Index int
}

type RoundingBehavior int

const (
	roundUp RoundingBehavior = iota
	roundDown
)

func New[T CmdMetadataAccessor](opts ...Option[T]) *Model[T] {
	textinput := textinput.New()
	textinput.Focus()
	formatters := DefaultFormatters()
	model := &Model[T]{
		textinput:          textinput,
		commandPlaceholder: []rune(""),
		subcommandWithArgs: "",
		prompt:             "> ",
		formatters:         formatters,
		parsedText:         &Statement{},
		delimiterRegex:     regexp.MustCompile(`\s+`),
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

func (m *Model[T]) SetFormatters(formatters Formatters) {
	m.formatters = formatters
}

func (m *Model[T]) NewPositionalArg(placeholder string) PositionalArg {
	return PositionalArg{
		placeholder:      placeholder,
		ArgStyle:         m.formatters.PositionalArg.Arg,
		PlaceholderStyle: m.formatters.PositionalArg.Placeholder,
	}
}

func (m *Model[T]) NewFlagPlaceholder(placeholder string) FlagPlaceholder {
	return FlagPlaceholder{
		text:  placeholder,
		Style: m.formatters.Flag.Placeholder,
	}
}

func (m *Model[T]) ShouldSelectSuggestion(suggestion input.Suggestion[T]) bool {
	currentTokenPos := m.CurrentTokenPos()
	currentToken := m.CurrentToken()
	// Only select if cursor is at the end of the token or the input will cut off the part after the cursor
	return m.CursorIndex() == currentTokenPos.End && currentToken == suggestion.Text
}

func (m *Model[T]) ShouldUnselectSuggestion(prevRunes []rune, msg tea.KeyMsg) bool {
	pos := m.CursorIndex()
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete:
		return pos < len(prevRunes) && !m.isDelimiter(string(prevRunes[pos]))
	case tea.KeyRunes, tea.KeySpace:
		if msg.String() != "=" {
			return true
		}
		token := ""
		if m.CursorIndex() == len(m.Runes()) {
			tokens := m.Tokens()
			token = tokens[len(tokens)-1].Value
		} else {
			token = m.CurrentTokenRoundDown()
		}
		// Don't unselect if the current token is a flag and we're adding an = delimiter
		return !strings.HasPrefix(token, "-")
	default:
		return true
	}
}

func (m *Model[T]) ShouldClearSuggestions(prevText []rune, msg tea.KeyMsg) bool {
	return m.isDelimiter(msg.String())
}

func (m *Model[T]) SelectedCommand() *input.Suggestion[T] {
	return m.selectedCommand
}

func (m *Model[T]) ArgsBeforeCursor() []string {
	args := []string{}
	runesBeforeCursor := m.Runes()[:m.CursorIndex()]

	expr, _ := m.parser.Parse(string(runesBeforeCursor))

	for _, arg := range expr.Args.Value {
		args = append(args, arg.Value)
	}
	return args
}

func (m *Model[T]) CompletedArgsBeforeCursor() []string {
	args := []string{}
	runesBeforeCursor := m.Runes()[:m.CursorIndex()]

	expr, _ := m.parser.Parse(string(runesBeforeCursor))

	for _, arg := range expr.Args.Value {
		args = append(args, arg.Value)
	}

	if len(expr.Flags.Value) == 0 && len(runesBeforeCursor) > 0 && !m.isDelimiter(string(runesBeforeCursor[len(runesBeforeCursor)-1])) {
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
	inputRunes := []rune(inputStr)
	suggestions := []input.Suggestion[T]{}
	isLong := strings.HasPrefix(inputStr, "--")
	isMulti := !isLong && strings.HasPrefix(inputStr, "-") && len(inputRunes) > 1
	tokenIndex := m.CurrentTokenPos().Index
	allTokens := m.Tokens()
	prevToken := ""
	if tokenIndex > 0 {
		prevToken = allTokens[tokenIndex-1].Value
	}

	currentIsFlag := false
	currentToken := ""
	if tokenIndex < len(allTokens) {
		currentToken = allTokens[tokenIndex].Value
		currentIsFlag = strings.HasPrefix(currentToken, "-")
	}

	curFlagText := ""
	if isMulti {
		curFlagText = string(inputRunes[len(inputRunes)-1])
	}

	for _, flag := range flags {
		// Don't show any flag suggestions if the current flag requires an arg unless the user skipped the arg and is now typing another flag that does not require an arg
		if ((isMulti && flag.Short == curFlagText) ||
			prevToken == "-"+flag.Short ||
			prevToken == "--"+flag.Long) && flag.RequiresArg() && (!currentIsFlag || currentToken == "-"+flag.Short || currentToken == "--"+flag.Long) {
			return []input.Suggestion[T]{}
		}

		long := "--" + flag.Long
		short := "-" + flag.Short
		if ((isLong || flag.Short == "") && strings.HasPrefix(long, inputStr)) ||
			strings.HasPrefix(short, inputStr) || (isMulti && !flag.RequiresArg()) {
			suggestion := input.Suggestion[T]{
				Description: flag.Description,
			}
			if isLong {
				suggestion.Text = long
			} else if isMulti {
				suggestion.Text = flag.Short
				// Ensure the completion text still has the leading dash for consistency
				suggestion.CompletionText = short
			} else {
				suggestion.Text = short
			}

			if suggestionFunc == nil {
				metadata := *new(T)
				placeholderField := reflect.ValueOf(&metadata).Elem().FieldByName("FlagPlaceholder")
				if placeholderField.IsValid() {
					placeholderField.Set(reflect.ValueOf(flag.Placeholder))
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
			text:             posArg.placeholder,
			placeholderStyle: posArg.PlaceholderStyle,
			argStyle:         posArg.ArgStyle,
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
			if m.currentFlag != nil && m.currentFlag.Metadata.GetFlagPlaceholder().text == "" {
				m.currentFlag = nil
			}
			// Clear any temporary placeholders
			m.args = m.origArgs
			return nil
		}

		if strings.HasPrefix(suggestion.Text, "-") {
			m.currentFlag = suggestion
		} else {
			m.showFlagPlaceholder = suggestion.Metadata.GetShowFlagPlaceholder()
			m.currentFlag = nil
		}
		index := m.CurrentTokenPos().Index

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
				placeholderStyle: m.formatters.Placeholder,
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
			m.commandPlaceholder = []rune("")
			m.subcommandWithArgs = ""
		} else {
			if !strings.HasPrefix(suggestion.Text, "-") {
				m.showFlagPlaceholder = suggestion.Metadata.GetShowFlagPlaceholder()
			}

			m.commandPlaceholder = []rune(suggestion.Text)
			m.subcommandWithArgs = suggestion.Text

			for _, posArg := range suggestion.Metadata.GetPositionalArgs() {
				newArg := arg{
					text:             posArg.placeholder,
					placeholderStyle: posArg.PlaceholderStyle,
					argStyle:         posArg.ArgStyle,
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
	token := m.CurrentToken()
	tokenRunes := []rune(token)
	suggestionRunes := []rune(suggestion.Text)
	tokenPos := m.CurrentTokenPos()

	if tokenPos.Index == 0 {
		m.selectedCommand = &suggestion
	}

	textRunes := m.Runes()
	if tokenPos.Start > -1 {
		cursor := m.CursorIndex()
		if strings.HasPrefix(token, "-") && strings.HasPrefix(suggestion.Text, "-") {
			// Adding an additional flag to the flag group, don't replace the entire token
			trailingRunes := []rune("")
			if cursor < len(textRunes) {
				// Add trailing text if we're not at the end of the line
				trailingRunes = textRunes[cursor+1:]
			}
			m.SetValue(string(textRunes[:cursor]) + suggestion.Text + string(trailingRunes))
		} else if strings.HasPrefix(token, "-") && !strings.HasPrefix(token, "--") && len(tokenRunes) > 2 && suggestion.Metadata.GetFlagPlaceholder().text == "" {
			// handle multi flag like -ab
			if cursor == tokenPos.Start {
				// If cursor is on the leading dash, replace the first two characters of the token ([-ab]c)
				m.SetValue(string(textRunes[:cursor]) + suggestion.Text + string(textRunes[cursor+2:]))
			} else {
				// If the cursor is after the dash, trim the dash from the suggestion and replace the single character on the cursor
				m.SetValue(string(textRunes[:cursor]) + string(suggestionRunes[1:]) + string(textRunes[cursor+1:]))
			}
		} else {
			m.SetValue(string(textRunes[:tokenPos.Start]) + suggestion.Text + string(textRunes[tokenPos.End:]))
			// Sometimes SetValue moves the cursor to the end of the line so we need to move it back to the current token
			m.SetCursor(len(textRunes[:tokenPos.Start]) + len(suggestionRunes) - suggestion.CursorOffset)
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

func (m *Model[T]) CompletionRunes(runes []rune) []rune {
	expr, _ := m.parser.Parse(string(runes))
	tokens := m.allTokens(expr)
	token := m.currentToken(tokens, roundUp)

	return token
}

func (m *Model[T]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *Model[T]) Value() string {
	return m.textinput.Value()
}

func (m *Model[T]) Runes() []rune {
	return []rune(m.textinput.Value())
}

func (m *Model[T]) ParsedValue() Statement {
	return *m.parsedText
}

func (m *Model[T]) CommandBeforeCursor() string {
	parsed := m.ParsedValue()
	commandRunes := []rune(parsed.Command.Value)
	if m.CursorIndex() >= len(commandRunes) {
		return parsed.Command.Value
	}
	return string(commandRunes[:m.CursorIndex()])
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

func (m Model[T]) Tokens() []input.Token {
	return m.allTokens(m.parsedText)
}

func (m Model[T]) AllTokensBeforeCursor() []input.Token {
	textBeforeCursor := m.Runes()[:m.CursorIndex()]

	expr, _ := m.parser.Parse(string(textBeforeCursor))
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

func (m Model[T]) allTokens(statement *Statement) []input.Token {
	tokens := []input.Token{statement.Command.ToToken(0, "command")}
	for i, arg := range statement.Args.Value {
		tokens = append(tokens, arg.ToToken(i+1, "arg"))
	}
	for _, flag := range statement.Flags.Value {
		tokens = append(tokens, input.TokenFromPos(flag.Name, "flag", len(tokens), flag.Pos))
		if flag.Value != nil {
			tokens = append(tokens, (*flag.Value).ToToken(len(tokens), "flagValue"))
		}
	}
	return tokens
}

func (m Model[T]) AllValues() []string {
	tokens := m.Tokens()
	values := []string{}
	for _, t := range tokens {
		values = append(values, t.Value)
	}
	return values
}

func (m Model[T]) CursorIndex() int {
	return m.textinput.Cursor()
}

func (m Model[T]) CursorOffset() int {
	cursorIndex := m.CursorIndex()
	runesBeforeCursor := m.Runes()[:cursorIndex]
	return runewidth.StringWidth(string(runesBeforeCursor))
}

func (m *Model[T]) SetCursor(pos int) {
	m.textinput.SetCursor(pos)
}

func (m Model[T]) Focused() bool {
	return m.textinput.Focused()
}

func (m *Model[T]) Prompt() string {
	return string(m.prompt)
}

func (m *Model[T]) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m Model[T]) cursorInToken(tokens []input.Token, pos int, roundingBehavior RoundingBehavior) bool {
	cursor := m.CursorIndex()
	isInToken := cursor >= tokens[pos].Start && cursor <= tokens[pos].End()
	if isInToken {
		return true
	}
	if roundingBehavior == roundDown {
		if pos == len(tokens)-1 {
			return true
		}
		return cursor < tokens[pos+1].Start
	} else {
		if pos == 0 {
			return false
		}
		return cursor > tokens[pos-1].End() && cursor < tokens[pos].Start
	}

}

func (m Model[T]) CurrentTokenPos() TokenPos {
	return m.currentTokenPos(m.Tokens(), roundUp)
}

func (m Model[T]) CurrentTokenPosRoundDown() TokenPos {
	return m.currentTokenPos(m.Tokens(), roundDown)
}

func (m Model[T]) currentTokenPos(tokens []input.Token, roundingBehavior RoundingBehavior) TokenPos {
	cursor := m.CursorIndex()
	if len(tokens) > 0 {
		last := tokens[len(tokens)-1]
		index := len(tokens) - 1
		runes := m.Runes()
		if roundingBehavior == roundUp && cursor > 0 && (m.isDelimiter(string(runes[cursor-1])) || (strings.HasPrefix(last.Value, "-") && string(runes[cursor-1]) == "=")) {
			// Haven't started a new token yet, but we have added a delimiter
			// so we'll consider the current token finished
			index++
		}
		// Check if cursor is at the end
		if cursor > last.End() {
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
				Start: tokens[i].Start,
				End:   tokens[i].End(),
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

func (m Model[T]) CurrentTokenBeforeCursor() string {
	return string(m.currentTokenBeforeCursor(roundUp))
}

func (m Model[T]) CurrentTokenBeforeCursorRoundDown() string {
	return string(m.currentTokenBeforeCursor(roundDown))
}

func (m Model[T]) currentTokenBeforeCursor(roundingBehavior RoundingBehavior) []rune {
	start := m.currentTokenPos(m.Tokens(), roundingBehavior).Start
	cursor := m.CursorIndex()
	if start > cursor {
		return []rune("")
	}
	val := m.Runes()[start:cursor]
	return val
}

func (m Model[T]) HasArgs() bool {
	return len(m.parsedText.Args.Value) > 0
}

func (m Model[T]) CurrentToken() string {
	return string(m.currentToken(m.Tokens(), roundUp))
}

func (m Model[T]) CurrentTokenRoundDown() string {
	return string(m.currentToken(m.Tokens(), roundDown))
}

func (m Model[T]) currentToken(tokens []input.Token, roundingBehavior RoundingBehavior) []rune {
	pos := m.currentTokenPos(tokens, roundingBehavior)
	return m.Runes()[pos.Start:pos.End]
}

func (m Model[T]) LastArg() *ident {
	parsed := *m.parsedText
	if len(parsed.Args.Value) == 0 {
		return nil
	}
	return &parsed.Args.Value[len(parsed.Args.Value)-1]
}

func (m Model[T]) CommandCompleted() bool {
	commandRunes := []rune(m.parsedText.Command.Value)
	if m.parsedText == nil || len(commandRunes) == 0 {
		return false
	}
	return m.CursorIndex() > m.parsedText.Command.Pos.Column-1+len(commandRunes)
}

func (m *Model[T]) Blur() {
	m.textinput.Blur()
}

func (m *Model[T]) OnExecutorFinished() {}

func (m Model[T]) View(viewMode input.ViewMode) string {
	viewBuilder := newCmdViewBuilder(m, viewMode)
	return viewBuilder.View()
}
