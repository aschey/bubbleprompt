package commandinput_test

import (
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/charmbracelet/lipgloss"
)

func ExampleDefaultFormatters() {
	defaultFormatters := commandinput.DefaultFormatters()
	defaultFormatters.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("128"))
}
