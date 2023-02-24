// Package commandinput provides an implementation of the [input.Input] interface.
// It should be used to build interactive CLI applications.
package commandinput

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// PositionalArg is a positional arg placeholder for completions.
type PositionalArg struct {
	placeholder string

	PlaceholderStyle lipgloss.Style
	ArgStyle         lipgloss.Style
}

// Placeholder returns the text value of the placeholder text.
func (p PositionalArg) Placeholder() string {
	return p.placeholder
}

// FlagArgPlaceholder is a flag placeholder for completions.
type FlagArgPlaceholder struct {
	text  string
	Style lipgloss.Style
}

// Text returns the placeholder text.
func (p FlagArgPlaceholder) Text() string {
	return p.text
}

// FlagInput is used to generate a list of flag suggestions.
type FlagInput struct {
	// Short is a short (single letter) flag with a single dash.
	// The leading dash can be optionally included.
	Short string
	// Long is a long (multi-letter) flag with multiple dashes.
	// The leading dashes can optionally be included.
	Long string
	// ArgPlaceholder is the placeholder for the flag argument (if applicable).
	ArgPlaceholder FlagArgPlaceholder
	// Description is the flag description.
	Description string
}

// ShortFlag returns the Short property formatted as a flag with a leading dash.
func (f FlagInput) ShortFlag() string {
	if len(f.Short) > 0 && !strings.HasPrefix(f.Short, "-") {
		return "-" + f.Short
	}
	return f.Short
}

// ShortFlag returns the Long property formatted as a flag with the leading dashes.
func (f FlagInput) LongFlag() string {
	if len(f.Long) > 0 && !strings.HasPrefix(f.Long, "--") {
		return "--" + f.Long
	}
	return f.Long
}

// RequiresArg returns whether or not the input has an argument placeholder.
// If no placeholder is supplied, then it is assumed that the [FlagInput] does not require an argument.
func (f FlagInput) RequiresArg() bool {
	return len(f.ArgPlaceholder.text) > 0
}

type modelState[T CommandMetadataAccessor] struct {
	selectedToken      *input.Token
	selectedSuggestion *suggestion.Suggestion[T]
	subcommand         *suggestion.Suggestion[T]
	selectedFlag       *suggestion.Suggestion[T]
	argNumber          int
}

func (m modelState[T]) isFlagSuggestion() bool {
	return (m.selectedSuggestion != nil && strings.HasPrefix(m.selectedSuggestion.Text, "-"))
}

func (m modelState[T]) isFlag() bool {
	return m.isFlagSuggestion() || m.selectedFlag != nil
}

// A Model is an input for handling CLI-style inputs.
// It provides advanced features such as placeholders and context-aware suggestions.
type Model[T CommandMetadataAccessor] struct {
	textinput        textinput.Model
	prompt           string
	defaultDelimiter string
	delimiterRegex   *regexp.Regexp
	formatters       Formatters
	parser           parser.Parser[statement]
	parsedText       *statement
	states           []modelState[T]
}

// New creates a new model.
func New[T CommandMetadataAccessor](opts ...Option[T]) *Model[T] {
	textinput := textinput.New()

	formatters := DefaultFormatters()
	model := &Model[T]{
		textinput:        textinput,
		prompt:           "> ",
		formatters:       formatters,
		parsedText:       &statement{},
		delimiterRegex:   regexp.MustCompile(`\s+`),
		defaultDelimiter: " ",
	}
	for _, opt := range opts {
		opt(model)
	}

	model.parser = buildCliParser(model.delimiterRegex.String())
	return model
}

// ParseUsage generates a list of [PositionalArg] from a usage string.
func (m *Model[T]) ParseUsage(placeholders string) ([]PositionalArg, error) {
	definition := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "RequiredArg", Pattern: `(<[^>]*>)`},
		{Name: "OptionalArg", Pattern: `\[[^\]]*\]`},
		{Name: "QuotedString", Pattern: `("[^"]*"?)|('[^']*'?)`},
		{Name: `String`, Pattern: `[^\s]+`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	lex, err := definition.LexString("", placeholders)
	if err != nil {
		return nil, err
	}
	tokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, err
	}
	positionalArgs := []PositionalArg{}
	for _, token := range tokens {
		if token.Type != lexer.EOF {
			positionalArgs = append(positionalArgs, m.NewPositionalArg(token.Value))
		}
	}

	return positionalArgs, nil
}

// Init is part of the [input.Input] interface. It should not be invoked by users of this library.
func (m *Model[T]) Init() tea.Cmd {
	return m.textinput.Focus()
}

// SetFormatters sets the formatters used by the input.
func (m *Model[T]) SetFormatters(formatters Formatters) {
	m.formatters = formatters
}

// Formatters returns the formatters used by the input.
func (m Model[T]) Formatters() Formatters {
	return m.formatters
}

// NewPositionalArg creates a positional arg placeholder for completions.
func (m *Model[T]) NewPositionalArg(placeholder string) PositionalArg {
	return PositionalArg{
		placeholder:      placeholder,
		ArgStyle:         m.formatters.PositionalArg.Arg,
		PlaceholderStyle: m.formatters.PositionalArg.Placeholder,
	}
}

// NewPositionalArgs creates multiple positional arg placeholders for completions.
func (m *Model[T]) NewPositionalArgs(placeholders ...string) []PositionalArg {
	args := []PositionalArg{}
	for _, placeholder := range placeholders {
		args = append(args, m.NewPositionalArg(placeholder))
	}
	return args
}

// NewFlagPlaceholder creates a flag placeholder for completions.
func (m *Model[T]) NewFlagPlaceholder(placeholder string) FlagArgPlaceholder {
	return FlagArgPlaceholder{
		text:  placeholder,
		Style: m.formatters.Flag.Placeholder,
	}
}

// ShouldSelectSuggestion is part of the [input.Input] interface.
// It should not be invoked by users of this library.
func (m *Model[T]) ShouldSelectSuggestion(suggestion suggestion.Suggestion[T]) bool {
	currentToken := m.CurrentToken()
	// Only select if cursor is at the end of the token or the input will cut off the part after the cursor
	return m.CursorIndex() == currentToken.End() && currentToken.Value == suggestion.Text
}

// ShouldUnselectSuggestion is part of the [input.Input] interface.
// It should not be invoked by users of this library.
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
			token = m.CurrentTokenRoundDown().Value
		}
		// Don't unselect if the current token is a flag and we're adding an = delimiter
		return !strings.HasPrefix(token, "-")
	default:
		return true
	}
}

// ShouldClearSuggestions is part of the [input.Input] interface.
// It should not be invoked by users of this library.
func (m *Model[T]) ShouldClearSuggestions(prevText []rune, msg tea.KeyMsg) bool {
	return m.isDelimiter(msg.String())
}

// ArgsBeforeCursor returns the positional arguments before the cursor position.
func (m *Model[T]) ArgsBeforeCursor() []string {
	args := []string{}
	runesBeforeCursor := m.Runes()[:m.CursorIndex()]

	expr, _ := m.parser.Parse(string(runesBeforeCursor))

	for _, arg := range expr.Args.Value {
		args = append(args, arg.Value)
	}
	return args
}

// CompletedArgsBeforeCursor returns the positional arguments before the cursor that have already been completed.
// In other words, there needs to be a delimiter after the argument to indicate that the user has finished
// entering in that argument.
func (m *Model[T]) CompletedArgsBeforeCursor() []string {
	args := []string{}
	runesBeforeCursor := m.Runes()[:m.CursorIndex()]

	expr, _ := m.parser.Parse(string(runesBeforeCursor))

	for _, arg := range expr.Args.Value {
		args = append(args, arg.Value)
	}

	if len(expr.Flags.Value) == 0 && len(runesBeforeCursor) > 0 &&
		!m.isDelimiter(string(runesBeforeCursor[len(runesBeforeCursor)-1])) {
		if len(args) > 0 {
			args = args[:len(args)-1]
		}
	}
	return args
}

// OnUpdateStart is part of the [input.Input] interface. It should not be invoked by users of this library.
func (m *Model[T]) OnUpdateStart(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	if _, ok := msg.(tea.KeyMsg); ok {
		expr, err := m.parser.Parse(m.Value())
		if err == nil {
			m.parsedText = expr
		}
	}
	allTokens := m.Tokens()
	tokenLen := len(allTokens)
	current := m.CurrentToken()
	if current.Value == "" {
		// Started a new token but haven't typed anything yet
		tokenLen += 1
	}

	if tokenLen < len(m.states) {
		m.states = m.states[:tokenLen]
	}
	if current.Index > len(m.states)-1 {
		newState := modelState[T]{}
		m.states = append(m.states, newState)
	}

	return cmd
}

// FlagSuggestions generates a list of [suggestion.Suggestion] based on
// the input string and the list of [FlagInput] supplied.
// The last parameter can be used to customize the metadata for the returned suggestions.
func (m *Model[T]) FlagSuggestions(
	inputStr string,
	flags []FlagInput,
	suggestionFunc func(FlagInput) T,
) []suggestion.Suggestion[T] {
	inputRunes := []rune(inputStr)
	suggestions := []suggestion.Suggestion[T]{}
	isLong := strings.HasPrefix(inputStr, "--")
	isMulti := !isLong && strings.HasPrefix(inputStr, "-") && len(inputRunes) > 1

	for _, flag := range flags {
		// Don't show any flag suggestions if the current flag requires an arg
		// unless the user skipped the arg and is now typing another flag that does not require an arg
		if m.shouldSkipFlagSuggestions(flag, inputRunes, isMulti) {
			return []suggestion.Suggestion[T]{}
		}

		if ((isLong || flag.Short == "") && strings.HasPrefix(flag.LongFlag(), inputStr)) ||
			strings.HasPrefix(flag.ShortFlag(), inputStr) || (isMulti && !flag.RequiresArg()) {

			suggestions = append(
				suggestions,
				m.getFlagSuggestion(flag, isLong, isMulti, suggestionFunc),
			)
		}
	}

	return suggestions
}

func (m *Model[T]) currentState() modelState[T] {
	index := m.CurrentToken().Index
	if index >= 0 {
		return m.states[index]
	} else {
		return modelState[T]{}
	}
}

func (m *Model[T]) shouldSkipFlagSuggestions(flag FlagInput, inputRunes []rune, isMulti bool) bool {
	tokenIndex := m.CurrentToken().Index
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
	return ((isMulti && flag.Short == curFlagText) ||
		prevToken == flag.ShortFlag() ||
		prevToken == flag.LongFlag()) && flag.RequiresArg() &&
		(!currentIsFlag || currentToken == flag.ShortFlag() || currentToken == flag.LongFlag())
}

func (m *Model[T]) getFlagSuggestion(
	flag FlagInput,
	isLong bool,
	isMulti bool,
	suggestionFunc func(FlagInput) T,
) suggestion.Suggestion[T] {
	suggestion := suggestion.Suggestion[T]{
		Description: flag.Description,
	}
	if isLong {
		suggestion.Text = flag.LongFlag()
	} else if isMulti {
		suggestion.Text = flag.Short
		// Ensure the suggestion text still has the leading dash for consistency
		suggestion.SuggestionText = flag.ShortFlag()
	} else {
		suggestion.Text = flag.ShortFlag()
	}

	if suggestionFunc == nil {
		metadata := *new(T)
		placeholderField := reflect.ValueOf(&metadata).Elem().FieldByName("FlagArgPlaceholder")
		if placeholderField.IsValid() {
			placeholderField.Set(reflect.ValueOf(flag.ArgPlaceholder))
			suggestion.Metadata = metadata
		}
	} else {
		suggestion.Metadata = suggestionFunc(flag)
	}

	return suggestion
}

// OnUpdateFinish is part of the [input.Input] interface.
// It should not be invoked by users of this library.
func (m *Model[T]) OnUpdateFinish(
	msg tea.Msg,
	suggestion *suggestion.Suggestion[T],
	isSelected bool,
) tea.Cmd {
	index := m.CurrentToken().Index

	m.states[index].selectedSuggestion = suggestion

	if index > 0 {
		subcommand := m.states[index-1].subcommand
		if subcommand != nil && m.states[index-1].argNumber+1 <= len(subcommand.Metadata.GetPositionalArgs()) {
			m.states[index].subcommand = m.states[index-1].subcommand
			m.states[index].argNumber = m.states[index-1].argNumber + 1
		}

		if m.states[index-1].isFlag() {
			m.states[index].selectedFlag = m.states[index-1].selectedFlag
		}
	}

	if suggestion != nil {
		if len(suggestion.Metadata.GetPositionalArgs()) > 0 {
			m.states[index].subcommand = suggestion
			m.states[index].argNumber = 0
		}
		if m.states[index].isFlagSuggestion() {
			m.states[index].selectedFlag = suggestion
		}
	}

	return nil
}

// OnSuggestionChanged is part of the [input.Input] interface. It should not be invoked by users of this library.
func (m *Model[T]) OnSuggestionChanged(suggestion suggestion.Suggestion[T]) {
	token := m.CurrentToken()
	tokenRunes := []rune(token.Value)
	suggestionRunes := []rune(suggestion.Text)
	m.states[token.Index].selectedToken = &token

	textRunes := m.Runes()
	if token.Index > -1 {
		cursor := m.CursorIndex()
		// Check if we're adding an additional flag to the flag group
		// If so, don't replace the entire token
		// Make sure the token already has at least one flag value appended to it first
		if strings.HasPrefix(token.Value, "-") &&
			!strings.HasPrefix(suggestion.Text, "-") {
			trailingRunes := []rune("")
			if cursor < len(textRunes) {
				// Add trailing text if we're not at the end of the line
				trailingRunes = textRunes[cursor+1:]
			}
			m.SetValue(string(textRunes[:cursor]) + suggestion.Text + string(trailingRunes))
		} else if strings.HasPrefix(token.Value, "-") &&
			!strings.HasPrefix(token.Value, "--") && len(tokenRunes) > 2 &&
			suggestion.Metadata.GetFlagArgPlaceholder().text == "" {
			// handle multi flag like -ab
			if cursor == token.Start {
				// If cursor is on the leading dash, replace the first two characters of the token ([-ab]c)
				m.SetValue(string(textRunes[:cursor]) + suggestion.Text + string(textRunes[cursor+2:]))
			} else {
				// If the cursor is after the dash, trim the dash from the suggestion and replace the single character on the cursor
				m.SetValue(string(textRunes[:cursor]) + string(suggestionRunes[1:]) + string(textRunes[cursor+1:]))
			}
		} else {
			m.SetValue(string(textRunes[:token.Start]) + suggestion.Text + string(textRunes[token.End():]))
			// Sometimes SetValue moves the cursor to the end of the line so we need to move it back to the current token
			m.SetCursor(len(textRunes[:token.Start]) + len(suggestionRunes) - suggestion.CursorOffset)
		}

	} else {
		m.SetValue(suggestion.Text)
	}
}

// OnSuggestionUnselected is part of the [input.Input] interface. It should not be invoked by users of this library.
func (m *Model[T]) OnSuggestionUnselected() {
	m.states[m.CurrentToken().Index].selectedToken = nil
}

// SuggestionRunes is part of the [input.Input] interface. It should not be invoked by users of this library.
func (m *Model[T]) SuggestionRunes(runes []rune) []rune {
	expr, _ := m.parser.Parse(string(runes))
	tokens := m.allTokens(expr)
	token := m.currentToken(tokens, input.RoundUp).Value

	return []rune(token)
}

// Focus sets the keyboard focus on the editor so the user can enter text.
func (m *Model[T]) Focus() tea.Cmd {
	return m.textinput.Focus()
}

// Focused returns whether the keyboard is focused on the input.
func (m Model[T]) Focused() bool {
	return m.textinput.Focused()
}

// Value returns the raw text entered by the user.
func (m *Model[T]) Value() string {
	return m.textinput.Value()
}

// Runes returns the raw text entered by the user as a list of runes.
// This is useful for indexing and length checks because doing these
// operations on strings does not work well with some unicode characters.
func (m *Model[T]) Runes() []rune {
	return []rune(m.textinput.Value())
}

// ParsedValue returns the input parsed into a [Statement].
func (m *Model[T]) ParsedValue() Statement {
	return (*m.parsedText).toStatement()
}

// CommandBeforeCursor returns the portion of the command (first input token) before the cursor position.
func (m *Model[T]) CommandBeforeCursor() string {
	parsed := m.parsedText
	commandRunes := []rune(parsed.Command.Value)
	if m.CursorIndex() >= len(commandRunes) {
		return parsed.Command.Value
	}
	return string(commandRunes[:m.CursorIndex()])
}

// SetValue overwrites the entire input with the given string.
func (m *Model[T]) SetValue(s string) {
	m.textinput.SetValue(s)
	expr, err := m.parser.Parse(m.Value())
	if err != nil {
		fmt.Println(err)
	}

	m.parsedText = expr
}

// ResetValue clears the entire input.
func (m *Model[T]) ResetValue() {
	m.textinput.SetValue("")
	m.parsedText = &statement{}
}

func (m *Model[T]) isDelimiter(s string) bool {
	return m.delimiterRegex.MatchString(s)
}

// Tokens returns the entire input as a list of [input.Token].
func (m Model[T]) Tokens() []input.Token {
	return m.allTokens(m.parsedText)
}

// TokensBeforeCursor returns the tokenized input before the cursor position.
func (m Model[T]) TokensBeforeCursor() []input.Token {
	textBeforeCursor := m.Runes()[:m.CursorIndex()]

	expr, _ := m.parser.Parse(string(textBeforeCursor))
	return m.allTokens(expr)
}

// ValuesBeforeCursor returns the token values of the entire input before the cursor position.
func (m Model[T]) ValuesBeforeCursor() []string {
	tokens := m.TokensBeforeCursor()
	values := []string{}
	for _, t := range tokens {
		values = append(values, t.Unquote())
	}
	return values
}

func (m Model[T]) allTokens(statement *statement) []input.Token {
	parsed := m.ParsedValue()
	tokens := []input.Token{parsed.Command}
	tokens = append(tokens, parsed.Args...)
	for _, flag := range parsed.Flags {
		tokens = append(tokens, flag.Name)
		if flag.Value != nil {
			tokens = append(tokens, *flag.Value)
		}
	}

	return tokens
}

// Values returns the tokenized input values.
func (m Model[T]) Values() []string {
	tokens := m.Tokens()
	values := []string{}
	for _, t := range tokens {
		values = append(values, t.Unquote())
	}
	return values
}

// CursorIndex returns the cursor index in terms of number of unicode characters.
// Use this to calculate input lengths in terms of number of characters entered.
func (m Model[T]) CursorIndex() int {
	return m.textinput.Position()
}

// CursorOffset returns the visual offset of the cursor in terms
// of number of terminal cells. Use this for calculating visual dimensions
// such as input width/height.
func (m Model[T]) CursorOffset() int {
	cursorIndex := m.CursorIndex()
	runesBeforeCursor := m.Runes()[:cursorIndex]
	return runewidth.StringWidth(string(runesBeforeCursor))
}

// SetCursor sets the cursor position.
func (m *Model[T]) SetCursor(pos int) {
	m.textinput.SetCursor(pos)
}

// SetCursorMode sets the mode of the cursor.
func (m *Model[T]) SetCursorMode(cursorMode cursor.Mode) tea.Cmd {
	return m.textinput.Cursor.SetMode(cursorMode)
}

// Prompt returns the terminal prompt.
func (m *Model[T]) Prompt() string {
	return string(m.prompt)
}

// SetPrompt sets the terminal prompt.
func (m *Model[T]) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m Model[T]) currentToken(
	tokens []input.Token,
	roundingBehavior input.RoundingBehavior,
) input.Token {
	return input.FindCurrentToken(
		m.Runes(),
		tokens,
		m.CursorIndex(),
		roundingBehavior,
		func(s string, last input.Token) bool {
			return m.isDelimiter(s) || (strings.HasPrefix(last.Value, "-") && s == "=")
		},
	)
}

// CurrentTokenBeforeCursor returns the portion of the current token before the cursor.
// If the cursor is between two tokens, it will take the token after the cursor.
func (m Model[T]) CurrentTokenBeforeCursor() input.Token {
	return m.currentTokenBeforeCursor(input.RoundUp)
}

// CurrentTokenBeforeCursorRoundDown returns the portion of the current token before the cursor.
// If the cursor is between two tokens, it will take the token before the cursor.
func (m Model[T]) CurrentTokenBeforeCursorRoundDown() input.Token {
	return m.currentTokenBeforeCursor(input.RoundDown)
}

// CurrentToken returns the token under the cursor.
// If the cursor is between two tokens, it will take the token after the cursor.
func (m Model[T]) CurrentToken() input.Token {
	return m.currentToken(m.Tokens(), input.RoundUp)
}

// CurrentTokenRoundDown returns the token under the cursor.
// If the cursor is between two tokens, it will take the token before the cursor.
func (m Model[T]) CurrentTokenRoundDown() input.Token {
	return m.currentToken(m.Tokens(), input.RoundDown)
}

func (m Model[T]) currentTokenBeforeCursor(roundingBehavior input.RoundingBehavior) input.Token {
	token := m.currentToken(m.Tokens(), roundingBehavior)
	start := token.Start
	cursor := m.CursorIndex()
	if start > cursor {
		return token
	}

	token.Value = string(m.Runes()[start:cursor])
	return token
}

// HasArgs returns whether the input has any positional arguments.
func (m Model[T]) HasArgs() bool {
	return len(m.parsedText.Args.Value) > 0
}

// LastArg returns the last positional argument in the input.
// If there are no arguments, it returns nil.
func (m Model[T]) LastArg() *input.Token {
	parsed := m.ParsedValue()
	if len(parsed.Args) == 0 {
		return nil
	}
	return &parsed.Args[len(parsed.Args)-1]
}

// CommandCompleted returns whether the user finished typing the entire command (first token).
func (m Model[T]) CommandCompleted() bool {
	commandRunes := []rune(m.parsedText.Command.Value)
	if m.parsedText == nil || len(commandRunes) == 0 {
		return false
	}
	return m.CursorIndex() > m.parsedText.Command.Pos.Column-1+len(commandRunes)
}

// Blur removes the focus from the input.
func (m *Model[T]) Blur() {
	m.textinput.Blur()
}

// OnExecutorFinished is part of the [input.Input] interface.
// It should not be invoked by users of this library.
func (m *Model[T]) OnExecutorFinished() {}

// View is part of the [input.Input] interface.
// It should not be invoked by users of this library.
func (m Model[T]) View(viewMode input.ViewMode) string {
	viewBuilder := newCmdViewBuilder(m, viewMode)
	return viewBuilder.View()
}
