package simpleinput_test

import "github.com/aschey/bubbleprompt/input/simpleinput"

func ExampleWithDelimiterRegex() {
	// Use period-delimited tokens instead of whitespace-delimited tokens
	simpleinput.New(
		simpleinput.WithDelimiterRegex[any](`\s*\.\s*`),
	)
}

func ExampleWithTokenRegex() {
	// Use period-delimited tokens instead of whitespace-delimited tokens
	simpleinput.New(
		simpleinput.WithTokenRegex[any](`[^\s\.]+`),
	)
}
