package commandinput_test

import (
	"fmt"

	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/charmbracelet/lipgloss"
)

func ExampleDefaultFormatters() {
	defaultFormatters := commandinput.DefaultFormatters()
	defaultFormatters.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("128"))
	fmt.Println(defaultFormatters.Cursor.GetForeground())

	// Output: 128
}
