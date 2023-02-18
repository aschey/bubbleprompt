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
	//b.renderFlagPlaceholder()
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
	args := append([]ident{b.model.parsedText.Command}, b.model.parsedText.Args.Value...)
	for i, arg := range args {

		if i < len(args) {
			if i == 0 {
				b.render([]rune(arg.Value), arg.Pos.Column, b.model.formatters.Command)
			} else {
				b.render([]rune(arg.Value), arg.Pos.Column, b.model.formatters.PositionalArg.Arg)
			}

		} else {
			state := b.model.states[i]
			b.renderCurrentArg(i, arg, state.selectedSuggestion)
		}

	}
}

func (b commandViewBuilder[T]) renderArg() {
}

func (b commandViewBuilder[T]) renderCurrentArg(argIndex int, arg ident, suggestion *suggestion.Suggestion[T]) {
	if suggestion != nil && strings.HasPrefix(suggestion.GetSuggestionText(), arg.Value) {
		tokenPos := len([]rune(arg.Value))
		suggestionRunes := []rune(suggestion.GetSuggestionText())
		b.render([]rune(suggestionRunes[tokenPos:]), b.viewBuilder.ViewLen(), b.model.formatters.Placeholder)
	}
}

func (b commandViewBuilder[T]) renderFlags() {
	flags := b.model.parsedText.Flags.Value
	currentFlagRunes := []rune{}
	currentFlagPlaceholderRunes := []rune{}
	currentState := b.model.currentState()
	if currentState.isFlagSuggestion() {
		currentFlagRunes = []rune(b.model.CurrentToken().Value)
		currentFlagPlaceholderRunes = []rune(
			currentState.selectedSuggestion.Metadata.GetFlagArgPlaceholder().text,
		)
	}

	for i, flag := range flags {
		b.renderFlag(i, flag, currentFlagRunes, currentFlagPlaceholderRunes)
	}
}

func (b commandViewBuilder[T]) renderFlag(
	i int,
	flag flag,
	currentFlagRunes []rune,
	currentFlagPlaceholderRunes []rune,
) {

	flagNameRunes := []rune(flag.Name)

	b.render(flagNameRunes, flag.Pos.Column, b.model.formatters.Flag.Flag)

	hasCurrentFlag := len(currentFlagRunes) > 0
	hasValue := flag.Value != nil
	// Render delimiter only once the full flag has been typed
	if !hasCurrentFlag || len(flagNameRunes) >= len(currentFlagRunes) || hasValue {
		if flag.Delim != nil {
			b.viewBuilder.Render(
				[]rune(flag.Delim.Value),
				flag.Delim.Pos.Column,
				lipgloss.NewStyle(),
			)
		}
	}

	b.renderFlagValue(i, flag, currentFlagRunes, currentFlagPlaceholderRunes)

}

func (b commandViewBuilder[T]) renderFlagValue(
	i int,
	flag flag,
	currentFlagRunes []rune,
	currentFlagPlaceholderRunes []rune,
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
		// Render current flag with placeholder info only if it's the last flag
		if len(currentFlagRunes) > 0 &&
			i == len(flags)-1 &&
			token.Start >= flag.Pos.Column-1 {
			b.renderLastFlag(flags, flag, currentFlagRunes, currentFlagPlaceholderRunes)
		}

		if flag.Value != nil {
			b.render(flagValueRunes, flag.Value.Pos.Column, b.flagValueStyle(flag.Value.Value))
		}
	}
}

func (b commandViewBuilder[T]) renderLastFlag(
	flags []flag,
	flag flag,
	currentFlagRunes []rune,
	currentFlagPlaceholderRunes []rune,
) {
	flagNameRunes := []rune(flag.Name)
	argVal := ""
	if len(flags) > 0 {
		argVal = flags[len(flags)-1].Name
	}

	// Render the rest of the arg placeholder only if the prefix matches
	if b.showPlaceholders && strings.HasPrefix(string(currentFlagRunes), argVal) {
		tokenPos := len(argVal)
		b.render(
			currentFlagRunes[tokenPos:],
			b.viewBuilder.ViewLen(),
			b.model.formatters.Placeholder,
		)
	}

	if len(currentFlagPlaceholderRunes) > 0 &&
		flagNameRunes[len(flagNameRunes)-1] != '-' {
		if !b.model.isDelimiter(string(*b.viewBuilder.Last())) && *b.viewBuilder.Last() != '=' {
			b.viewBuilder.Render(
				[]rune(b.model.defaultDelimiter),
				b.viewBuilder.ViewLen(),
				lipgloss.NewStyle(),
			)
		}

		if b.showPlaceholders && flag.Value == nil {
			b.viewBuilder.RenderPlaceholder(
				currentFlagPlaceholderRunes,
				b.viewBuilder.ViewLen(),
				b.model.formatters.Flag.Placeholder,
			)
		}
	}
}

// func (b commandViewBuilder[T]) renderFlagPlaceholder() {
// 	if b.showPlaceholders && len(b.model.parsedText.Flags.Value) == 0 &&
// 		b.currentState.selectedSuggestion.Metadata.GetShowFlagPlaceholder() {
// 		b.viewBuilder.Render([]rune("[flags]"), b.viewBuilder.ViewLen(), b.model.formatters.Placeholder)
// 		return
// 	}
// 	if b.model.CurrentToken().Index < len(b.model.Tokens())-1 {
// 		return
// 	}
// 	subcommand := b.model.currentState().subcommand
// 	if subcommand == nil {
// 		return
// 	}
// 	if !b.model.currentState().isFlagSuggestion() {
// 		return
// 	}
// 	flagSuggestion := b.model.currentState().selectedSuggestion.Metadata.GetFlagArgPlaceholder()
// 	if flagSuggestion.text == "" {
// 		return
// 	}
// 	if strings.HasPrefix(string(subcommand.GetSuggestionText()), b.model.Tokens()[b.model.currentState().argNumber].Value) {
// 		last := b.viewBuilder.Last()
// 		if last == nil || !b.model.isDelimiter(string(*last)) {
// 			b.viewBuilder.Render(
// 				[]rune(b.model.defaultDelimiter),
// 				b.viewBuilder.ViewLen(),
// 				lipgloss.NewStyle(),
// 			)
// 		}
// 		b.viewBuilder.Render(
// 			[]rune(flagSuggestion.text),
// 			b.viewBuilder.ViewLen(),
// 			b.model.Formatters().Placeholder,
// 		)
// 	}
// }

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
