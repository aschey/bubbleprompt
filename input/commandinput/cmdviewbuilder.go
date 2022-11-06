package commandinput

import (
	"strings"

	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/lipgloss"
)

type cmdViewBuilder[T CmdMetadataAccessor] struct {
	model            Model[T]
	viewBuilder      *input.ViewBuilder
	showPlaceholders bool
	showCursor       bool
}

func newCmdViewBuilder[T CmdMetadataAccessor](model Model[T], viewMode input.ViewMode) cmdViewBuilder[T] {
	showCursor := !model.textinput.Blink()
	if viewMode == input.Static {
		showCursor = false
	}
	showPlaceholders := viewMode == input.Interactive
	viewBuilder := input.NewViewBuilder(model.Cursor(), model.CursorStyle, model.defaultDelimiter, showCursor)
	return cmdViewBuilder[T]{
		model, viewBuilder, showPlaceholders, showCursor,
	}
}

func (b cmdViewBuilder[T]) View() string {
	b.renderCommand()
	b.renderPrefix()
	b.renderArgs()
	b.renderCurrentArg()
	b.renderFlags()
	b.renderFlagPlaceholders()
	b.renderTrailingText()

	return b.model.PromptStyle.Render(b.model.prompt) + b.viewBuilder.View()
}

func (b cmdViewBuilder[T]) renderCommand() {
	command := b.model.parsedText.Command.Value
	if b.model.selectedCommand == nil {
		b.viewBuilder.Render(command, b.model.parsedText.Command.Pos.Offset, b.model.TextStyle)
	} else {
		b.viewBuilder.Render(command, b.model.parsedText.Command.Pos.Offset, b.model.SelectedTextStyle)
	}
}

func (b cmdViewBuilder[T]) renderPrefix() {
	command := b.model.parsedText.Command.Value
	if b.showPlaceholders && strings.HasPrefix(b.model.commandPlaceholder, b.model.Value()) && b.model.commandPlaceholder != command {
		b.viewBuilder.Render(b.model.commandPlaceholder[len(command):], b.model.parsedText.Command.Pos.Offset+len(command), b.model.PlaceholderStyle)
	}
}

func (b cmdViewBuilder[T]) renderArgs() {
	for i, arg := range b.model.parsedText.Args.Value {
		argStyle := lipgloss.NewStyle()
		if i < len(b.model.args) {
			argStyle = b.model.args[i].argStyle
		}
		b.viewBuilder.Render(arg.Value, arg.Pos.Offset, argStyle)
	}
}

func (b cmdViewBuilder[T]) renderCurrentArg() {
	// Render current arg if persist == true
	currentArg := len(b.model.parsedText.Args.Value) - 1
	if currentArg >= 0 && currentArg < len(b.model.args) {
		arg := b.model.args[currentArg]
		argVal := b.model.parsedText.Args.Value[currentArg].Value
		// Render the rest of the arg placeholder only if the prefix matches
		if arg.persist && strings.HasPrefix(arg.text, argVal) {
			tokenPos := len(argVal)
			b.viewBuilder.Render(arg.text[tokenPos:], b.viewBuilder.ViewLen(), arg.placeholderStyle)
		}
	}
}

func (b cmdViewBuilder[T]) renderFlags() {
	currentPos := b.model.CurrentTokenPos(RoundDown).Start
	currentToken := b.model.CurrentToken(RoundDown)
	for i, flag := range b.model.parsedText.Flags.Value {
		b.viewBuilder.Render(flag.Name, flag.Pos.Offset, lipgloss.NewStyle().Foreground(lipgloss.Color("245")))
		// Render delimiter only once the full flag has been typed
		if b.model.currentFlag == nil || len(flag.Name) >= len(b.model.currentFlag.Text) || flag.Value != nil {
			if flag.Delim != nil {
				b.viewBuilder.Render(flag.Delim.Value, flag.Delim.Pos.Offset, lipgloss.NewStyle())
			}
		}

		if (flag.Pos.Offset != currentPos) && (currentPos < b.model.Cursor() || i < len(b.model.parsedText.Flags.Value)-1 || (len(currentToken) > 0 && !strings.HasPrefix(currentToken, "-"))) {
			if flag.Value != nil {
				b.viewBuilder.Render(flag.Value.Value, flag.Value.Pos.Offset, b.model.flagValueStyle(flag.Value.Value))
			}

		} else {
			// Render current flag
			if b.model.currentFlag != nil {
				argVal := ""
				if len(b.model.parsedText.Flags.Value) > 0 {
					argVal = b.model.parsedText.Flags.Value[len(b.model.parsedText.Flags.Value)-1].Name
				}

				// Render the rest of the arg placeholder only if the prefix matches
				if b.showPlaceholders && strings.HasPrefix(b.model.currentFlag.Text, argVal) {
					tokenPos := len(argVal)
					b.viewBuilder.Render(b.model.currentFlag.Text[tokenPos:], b.viewBuilder.ViewLen(), b.model.PlaceholderStyle)
				}
			}

			if b.model.currentFlag != nil && len(b.model.currentFlag.Metadata.GetFlagPlaceholder().Text) > 0 && flag.Name[len(flag.Name)-1] != '-' {
				if !b.model.isDelimiter(string(*b.viewBuilder.Last())) && *b.viewBuilder.Last() != '=' {
					b.viewBuilder.Render(b.model.defaultDelimiter, b.viewBuilder.ViewLen(), lipgloss.NewStyle())
				}

				if b.showPlaceholders && flag.Value == nil {
					b.viewBuilder.RenderPlaceholder(b.model.currentFlag.Metadata.GetFlagPlaceholder().Text, b.viewBuilder.ViewLen(), lipgloss.NewStyle().Foreground(lipgloss.Color("14")))
				}

			}
			if flag.Value != nil {
				b.viewBuilder.Render(flag.Value.Value, flag.Value.Pos.Offset, b.model.flagValueStyle(flag.Value.Value))
			}
		}
	}
}

func (b cmdViewBuilder[T]) renderFlagPlaceholders() {
	args := b.model.args
	if b.showPlaceholders && len(b.model.parsedText.Flags.Value) == 0 && b.model.showFlagPlaceholder {
		args = append(args, arg{text: "[flags]", placeholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("14"))})
	}

	// Render arg placeholders
	startPlaceholder := len(b.model.parsedText.Args.Value) + len(b.model.parsedText.Flags.Value)
	// Don't show arg placeholders if the current arg doesn't match the arg we're about to show placeholders for (the user moved the cursor over to the left)
	all := b.model.AllValues()
	if b.model.suggestionLevel > len(all)-1 || strings.HasPrefix(b.model.subcommandWithArgs, all[b.model.suggestionLevel]) {
		if b.showPlaceholders && startPlaceholder < len(args) {
			for _, arg := range args[startPlaceholder:] {
				last := b.viewBuilder.Last()
				if last == nil || !b.model.isDelimiter(string(*last)) {
					b.viewBuilder.Render(b.model.defaultDelimiter, b.viewBuilder.ViewLen(), lipgloss.NewStyle())
				}

				b.viewBuilder.Render(arg.text, b.viewBuilder.ViewLen(), arg.placeholderStyle)
			}
		}
	}
}

func (b cmdViewBuilder[T]) renderTrailingText() {
	value := b.model.Value()
	viewLen := b.viewBuilder.ViewLen()
	if len(value) > viewLen {
		b.viewBuilder.Render(value[b.viewBuilder.ViewLen():], b.viewBuilder.ViewLen(), lipgloss.NewStyle())
	}

}
