package prompt

import (
	"github.com/aschey/bubbleprompt/commandinput"
)

type Document struct {
	Text           string
	ParsedInput    commandinput.Statement
	CursorPosition int
}

func (d Document) TextBeforeCursor() string {
	if d.CursorPosition >= len(d.Text) {
		return d.Text
	}
	return d.Text[:d.CursorPosition]
}

func (d Document) CommandBeforeCursor() string {
	if d.CursorPosition >= len(d.ParsedInput.Command.Value) {
		return d.ParsedInput.Command.Value
	}
	return d.ParsedInput.Command.Value[:d.CursorPosition]
}
