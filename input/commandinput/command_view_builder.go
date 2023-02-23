package commandinput

import (
	"strconv"
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/lipgloss"
)

type commandViewBuilder[T CommandMetadataAccessor] struct {
	model            Model[T]
	viewBuilder      *input.ViewBuilder
	currentState     modelState[T]
	showPlaceholders bool
	showCursor       bool
}

func newCmdViewBuilder[T CommandMetadataAccessor](
	model Model[T],
	viewMode input.ViewMode,
) commandViewBuilder[T] {
	showCursor := !model.textinput.Cursor.Blink
	if viewMode == input.Static {
		showCursor = false
	}
	showPlaceholders := viewMode == input.Interactive
	viewBuilder := input.NewViewBuilder(
		model.CursorIndex(),
		model.formatters.Cursor,
		model.defaultDelimiter,
		showCursor,
	)
	return commandViewBuilder[T]{
		model, viewBuilder, model.currentState(), showPlaceholders, showCursor,
	}
}

func (b commandViewBuilder[T]) View() string {
	b.renderArgs()
	b.renderFlags()
	b.renderFlagsPlaceholder()
	b.renderPlaceholders()
	b.renderFlagPlaceholder()
	b.renderTrailingText()

	return b.model.formatters.Prompt.Render(string(b.model.prompt)) + b.viewBuilder.View()
}

func (b commandViewBuilder[T]) render(runes []rune, column int, style lipgloss.Style) {
	if b.currentState.selectedToken != nil && b.currentState.selectedToken.Start == column-1 {
		b.viewBuilder.Render(runes, column, b.model.formatters.SelectedText)
	} else {
		b.viewBuilder.Render(runes, column, style)
	}
}

func (b commandViewBuilder[T]) renderArgs() {
	command := b.model.parsedText.Command
	args := b.model.parsedText.Args.Value
	b.render([]rune(command.Value), command.Pos.Column, b.model.formatters.Command)
	if len(args) == 0 {
		b.renderCurrentArg(command.Value, b.model.states[0].selectedSuggestion)
	} else {
		for _, arg := range args {
			b.render([]rune(arg.Value), arg.Pos.Column, b.model.formatters.PositionalArg.Arg)
		}

		b.renderCurrentArg(args[len(args)-1].Value, b.model.states[len(args)].selectedSuggestion)
	}
}

func (b commandViewBuilder[T]) renderCurrentArg(arg string, suggestion *suggestion.Suggestion[T]) {
	if len(arg) > 0 && suggestion != nil && strings.HasPrefix(suggestion.GetSuggestionText(), arg) {
		tokenPos := len([]rune(arg))
		suggestionRunes := []rune(suggestion.GetSuggestionText())
		b.render([]rune(suggestionRunes[tokenPos:]), b.viewBuilder.ViewLen(), b.model.formatters.Placeholder)
	}
}

func (b commandViewBuilder[T]) renderFlags() {
	flags := b.model.parsedText.Flags.Value

	for i, flag := range flags {
		b.renderFlag(i, flag)
	}

	if len(flags) > 0 && b.currentState.isFlagSuggestion() {
		b.renderCurrentArg(flags[len(flags)-1].Name, b.currentState.selectedSuggestion)
	}
}

func (b commandViewBuilder[T]) renderFlag(
	i int,
	flag flag,
) {
	flagNameRunes := []rune(flag.Name)

	b.render(flagNameRunes, flag.Pos.Column, b.model.formatters.Flag.Flag)

	hasValue := flag.Value != nil
	// Render delimiter only once the full flag has been typed
	if hasValue {
		if flag.Delim != nil {
			b.viewBuilder.Render(
				[]rune(flag.Delim.Value),
				flag.Delim.Pos.Column,
				lipgloss.NewStyle(),
			)
		}
	}

	b.renderFlagValue(i, flag)
}

func (b commandViewBuilder[T]) renderFlagValue(
	i int,
	flag flag,
) {
	token := b.model.CurrentTokenRoundDown()
	flags := b.model.parsedText.Flags.Value
	flagIsCurrent := flag.Pos.Column-1 == token.Start
	cursorBeforeEnd := token.Start < b.model.CursorIndex()
	beforeLastFlag := i < len(flags)-1
	currentTokenIsNotFlag := len(token.Value) > 0 && !strings.HasPrefix(token.Value, "-")
	flagValueRunes := []rune{}
	if flag.Value != nil {
		flagValueRunes = []rune(flag.Value.Value)
	}

	if !flagIsCurrent && (cursorBeforeEnd || beforeLastFlag || currentTokenIsNotFlag) {
		if flag.Value != nil {
			b.render(flagValueRunes, flag.Value.Pos.Column, b.flagValueStyle(flag.Value.Value))
		}
	} else {
		if flag.Value != nil {
			b.render(flagValueRunes, flag.Value.Pos.Column, b.flagValueStyle(flag.Value.Value))
		}
	}
}

func (b commandViewBuilder[T]) renderDelimiter() {
	last := b.viewBuilder.Last()
	if last != nil && !b.model.isDelimiter(string(*last)) {
		b.viewBuilder.Render(
			[]rune(b.model.defaultDelimiter),
			b.viewBuilder.ViewLen(),
			lipgloss.NewStyle(),
		)
	}
}

func (b commandViewBuilder[T]) renderFlagDelimiter() {
	last := b.viewBuilder.Last()
	if last != nil && !(b.model.isDelimiter(string(*last)) || *last == '=') {
		b.viewBuilder.Render(
			[]rune(b.model.defaultDelimiter),
			b.viewBuilder.ViewLen(),
			lipgloss.NewStyle(),
		)
	}
}

func (b commandViewBuilder[T]) renderFlagsPlaceholder() {
	if b.showPlaceholders && len(b.model.parsedText.Flags.Value) == 0 && b.currentState.selectedSuggestion != nil &&
		b.currentState.selectedSuggestion.Metadata.GetShowFlagPlaceholder() {
		b.renderDelimiter()
		b.viewBuilder.Render([]rune("[flags]"), b.viewBuilder.ViewLen(), b.model.formatters.Flag.Placeholder)
	}
}

func (b commandViewBuilder[T]) renderPlaceholders() {
	currentToken := b.model.CurrentToken()
	tokenLen := len(b.model.Tokens())
	if currentToken.Index < tokenLen-1 {
		return
	}

	currentState := b.model.currentState()
	if currentToken.Value == "" && currentState.selectedSuggestion != nil {
		b.renderDelimiter()
		b.viewBuilder.Render(
			[]rune(currentState.selectedSuggestion.GetSuggestionText()),
			b.viewBuilder.ViewLen(),
			b.model.formatters.Placeholder)
	}

	if currentState.subcommand == nil {
		return
	}

	positionalArgs := currentState.subcommand.Metadata.GetPositionalArgs()
	argNumber := currentState.argNumber
	if len(currentToken.Value) == 0 && argNumber > 0 {
		argNumber--
	}
	if argNumber < len(positionalArgs) {
		for _, arg := range positionalArgs[argNumber:] {
			b.renderDelimiter()
			b.viewBuilder.Render(
				[]rune(arg.placeholder),
				b.viewBuilder.ViewLen(),
				arg.PlaceholderStyle,
			)
		}
	}
}

func (b commandViewBuilder[T]) renderFlagPlaceholder() {
	currentState := b.model.currentState()
	if currentState.isFlagSuggestion() {
		flagArgPlaceholder := currentState.selectedSuggestion.Metadata.GetFlagArgPlaceholder()
		if flagArgPlaceholder.text != "" {
			b.renderFlagDelimiter()
			b.viewBuilder.Render(
				[]rune(flagArgPlaceholder.text),
				b.viewBuilder.ViewLen(),
				flagArgPlaceholder.Style)

		}
	}
}

func (b commandViewBuilder[T]) renderTrailingText() {
	value := []rune(b.model.Value())
	viewLen := b.viewBuilder.ViewLen()
	if len(value) > viewLen {
		b.viewBuilder.Render(
			value[b.viewBuilder.ViewLen():],
			b.viewBuilder.ViewLen(),
			lipgloss.NewStyle(),
		)
	}
}

func (b commandViewBuilder[T]) flagValueStyle(value string) lipgloss.Style {
	if _, err := strconv.ParseInt(value, 10, 32); err == nil {
		return b.model.formatters.FlagValue.Number
	} else if _, err := strconv.ParseBool(value); err == nil {
		return b.model.formatters.FlagValue.Bool
	}
	return b.model.formatters.FlagValue.String
}
