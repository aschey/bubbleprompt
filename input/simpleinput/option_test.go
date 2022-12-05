package simpleinput_test

import (
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/charmbracelet/lipgloss"
)

func ExampleWithDelimiterRegex() {
	// Use period-delimited tokens instead of whitespace-delimited tokens
	// If you change the delimiter regex, you'll probably also need to change the token regex
	simpleinput.New(
		simpleinput.WithTokenRegex[any](`[^\s\.]+`),
		simpleinput.WithDelimiterRegex[any](`\s*\.\s*`),
	)
}

func ExampleWithTokenRegex() {
	// Use period-delimited tokens instead of whitespace-delimited tokens
	// If you change the token regex, you'll probably also need to change the delimiter regex
	simpleinput.New(
		simpleinput.WithTokenRegex[any](`[^\s\.]+`),
		simpleinput.WithDelimiterRegex[any](`\s*\.\s*`),
	)
}

func ExampleWithSelectedTextStyle() {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	simpleinput.New(simpleinput.WithSelectedTextStyle[any](style))
}
