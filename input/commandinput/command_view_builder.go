package commandinput

import (
	"strconv"
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/lipgloss"
)

type commandViewBuilder[T CommandMetadataAccessor] struct {
	model            Model[T]
	viewBuilder      *input.ViewBuilder
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
		model, viewBuilder, showPlaceholders, showCursor,
	}
}

func (b commandViewBuilder[T]) View() string {
	b.renderCommand()
	b.renderPrefix()
	b.renderArgs()
	b.renderCurrentArg()
	b.renderFlags()
	b.renderFlagPlaceholders()
	b.renderTrailingText()

	return b.model.formatters.Prompt.Render(string(b.model.prompt)) + b.viewBuilder.View()
}

func (b commandViewBuilder[T]) render(runes []rune, column int, style lipgloss.Style) {
	if b.model.selectedToken != nil && b.model.selectedToken.Start == column-1 {
		b.viewBuilder.Render(runes, column, b.model.formatters.SelectedText)
	} else {
		b.viewBuilder.Render(runes, column, style)
	}
}

func (b commandViewBuilder[T]) renderCommand() {
	commandRunes := []rune(b.model.parsedText.Command.Value)
	b.render(commandRunes, b.model.parsedText.Command.Pos.Column, b.model.formatters.Command)
}

func (b commandViewBuilder[T]) renderPrefix() {
	commandRunes := []rune(b.model.parsedText.Command.Value)
	if b.showPlaceholders &&
		strings.HasPrefix(string(b.model.commandPlaceholder), b.model.Value()) &&
		string(b.model.commandPlaceholder) != string(commandRunes) {
		b.viewBuilder.Render(
			b.model.commandPlaceholder[len(commandRunes):],
			b.model.parsedText.Command.Pos.Column+len(commandRunes),
			b.model.formatters.Placeholder,
		)
	}
}

func (b commandViewBuilder[T]) renderArgs() {
	for i, arg := range b.model.parsedText.Args.Value {
		argStyle := lipgloss.NewStyle()
		if i < len(b.model.args) {
			argStyle = b.model.args[i].argStyle
		}
		b.render([]rune(arg.Value), arg.Pos.Column, argStyle)
	}
}

func (b commandViewBuilder[T]) renderCurrentArg() {
	// Render current arg if persist == true
	currentArg := len(b.model.parsedText.Args.Value) - 1
	if currentArg >= 0 && currentArg < len(b.model.args) {
		arg := b.model.args[currentArg]
		argVal := b.model.parsedText.Args.Value[currentArg].Value
		// Render the rest of the arg placeholder only if the prefix matches
		if arg.persist && strings.HasPrefix(arg.text, argVal) {
			tokenPos := len([]rune(argVal))
			b.render([]rune(arg.text)[tokenPos:], b.viewBuilder.ViewLen(), arg.placeholderStyle)
		}
	}
}

func (b commandViewBuilder[T]) renderFlags() {
	flags := b.model.parsedText.Flags.Value
	currentFlagRunes := []rune{}
	currentFlagPlaceholderRunes := []rune{}
	if b.model.currentFlag != nil {
		currentFlagRunes = []rune(b.model.currentFlag.Text)
		currentFlagPlaceholderRunes = []rune(
			b.model.currentFlag.Metadata.GetFlagArgPlaceholder().text,
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

	hasCurrentFlag := b.model.currentFlag != nil
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
		if b.model.currentFlag != nil &&
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
	if b.showPlaceholders && strings.HasPrefix(b.model.currentFlag.Text, argVal) {
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

func (b commandViewBuilder[T]) renderFlagPlaceholders() {
	args := b.model.args
	if b.showPlaceholders && len(b.model.parsedText.Flags.Value) == 0 &&
		b.model.showFlagPlaceholder {
		args = append(
			args,
			arg{text: "[flags]", placeholderStyle: b.model.formatters.Flag.Placeholder},
		)
	}

	// Render arg placeholders
	startPlaceholder := len(b.model.parsedText.Args.Value) + len(b.model.parsedText.Flags.Value)
	// Don't show arg placeholders if the current arg doesn't match the arg
	// we're about to show placeholders for (the user moved the cursor over to the left)
	all := b.model.Values()
	if b.model.suggestionLevel > len(all)-1 ||
		strings.HasPrefix(string(b.model.subcommandWithArgs), all[b.model.suggestionLevel]) {
		if b.showPlaceholders && startPlaceholder < len(args) {
			for _, arg := range args[startPlaceholder:] {
				last := b.viewBuilder.Last()
				if last == nil || !b.model.isDelimiter(string(*last)) {
					b.viewBuilder.Render(
						[]rune(b.model.defaultDelimiter),
						b.viewBuilder.ViewLen(),
						lipgloss.NewStyle(),
					)
				}

				b.viewBuilder.Render(
					[]rune(arg.text),
					b.viewBuilder.ViewLen(),
					arg.placeholderStyle,
				)
			}
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
