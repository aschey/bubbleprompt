package commandinput_test

import (
	"fmt"

	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
)

func ExampleModel_NewPositionalArg() {
	textInput := commandinput.New[commandinput.CommandMetadata]()
	commandMetadata := commandinput.CommandMetadata{
		PositionalArgs: []commandinput.PositionalArg{textInput.NewPositionalArg("<arg1>")},
	}

	suggestions := []suggestion.Suggestion[commandinput.CommandMetadata]{
		{Text: "test", Metadata: commandMetadata},
	}

	fmt.Println(suggestions[0].Metadata.GetPositionalArgs()[0].Placeholder())
	// Output: <arg1>
}

func ExampleModel_NewPositionalArgs() {
	textInput := commandinput.New[commandinput.CommandMetadata]()
	commandMetadata := commandinput.CommandMetadata{
		PositionalArgs: textInput.NewPositionalArgs("<arg1>", "<arg2>"),
	}

	suggestions := []suggestion.Suggestion[commandinput.CommandMetadata]{
		{Text: "test", Metadata: commandMetadata},
	}

	fmt.Println(suggestions[0].Metadata.GetPositionalArgs()[0].Placeholder())
	fmt.Println(suggestions[0].Metadata.GetPositionalArgs()[1].Placeholder())
	// Output:
	// <arg1>
	// <arg2>
}

func ExampleModel_NewFlagPlaceholder() {
	textInput := commandinput.New[commandinput.CommandMetadata]()

	flags := []commandinput.FlagInput{
		{
			Short:       "d",
			Long:        "days",
			Placeholder: textInput.NewFlagPlaceholder("<number of days>"),
			Description: "Forecast days",
		},
	}

	fmt.Println(flags[0].Placeholder.Text())
	// Output: <number of days>
}

func ExampleModel_FlagSuggestions() {
	textInput := commandinput.New[commandinput.CommandMetadata]()
	flags := []commandinput.FlagInput{
		{
			Short:       "i",
			Long:        "interval",
			Description: "refresh interval",
			Placeholder: textInput.NewFlagPlaceholder("<value>"),
		},
	}

	suggestions := textInput.FlagSuggestions("", flags, nil)
	fmt.Printf("Text: %s, Description: %s\n", suggestions[0].Text, suggestions[0].Description)

	suggestions = textInput.FlagSuggestions("--", flags, nil)
	fmt.Printf("Text: %s, Description: %s\n", suggestions[0].Text, suggestions[0].Description)

	suggestions = textInput.FlagSuggestions("", flags,
		func(flagInput commandinput.FlagInput) commandinput.CommandMetadata {
			return commandinput.CommandMetadata{
				FlagPlaceholder:     flagInput.Placeholder,
				PreservePlaceholder: true,
			}
		})
	fmt.Printf("Text: %s, Description: %s, Preserve Placeholder: %t\n",
		suggestions[0].Text, suggestions[0].Description, suggestions[0].Metadata.PreservePlaceholder)

	// Output:
	// Text: -i, Description: refresh interval
	// Text: --interval, Description: refresh interval
	// Text: -i, Description: refresh interval, Preserve Placeholder: true
}
