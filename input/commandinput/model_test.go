package commandinput_test

import (
	"fmt"

	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
)

func ExampleModel_NewPositionalArg() {
	textInput := commandinput.New[commandinput.CommandMetadata[any]]()
	commandMetadata := commandinput.CommandMetadata[any]{
		PositionalArgs: []commandinput.PositionalArg{textInput.NewPositionalArg("<arg1>")},
	}

	suggestions := []suggestion.Suggestion[commandinput.CommandMetadata[any]]{
		{Text: "test", Metadata: commandMetadata},
	}

	fmt.Println(suggestions[0].Metadata.PositionalArgs[0].Placeholder())
	// Output: <arg1>
}

func ExampleModel_NewPositionalArgs() {
	textInput := commandinput.New[commandinput.CommandMetadata[any]]()
	commandMetadata := commandinput.CommandMetadata[any]{
		PositionalArgs: textInput.NewPositionalArgs("<arg1>", "<arg2>"),
	}

	suggestions := []suggestion.Suggestion[commandinput.CommandMetadata[any]]{
		{Text: "test", Metadata: commandMetadata},
	}

	fmt.Println(suggestions[0].Metadata.PositionalArgs[0].Placeholder())
	fmt.Println(suggestions[0].Metadata.PositionalArgs[1].Placeholder())
	// Output:
	// <arg1>
	// <arg2>
}

func ExampleModel_NewFlagPlaceholder() {
	textInput := commandinput.New[commandinput.CommandMetadata[any]]()

	flags := []commandinput.FlagInput{
		{
			Short:          "d",
			Long:           "days",
			ArgPlaceholder: textInput.NewFlagPlaceholder("<number of days>"),
			Description:    "Forecast days",
		},
	}

	fmt.Println(flags[0].ArgPlaceholder.Text())
	// Output: <number of days>
}

func ExampleModel_FlagSuggestions() {
	textInput := commandinput.New[any]()
	flags := []commandinput.FlagInput{
		{
			Short:          "i",
			Long:           "interval",
			Description:    "refresh interval",
			ArgPlaceholder: textInput.NewFlagPlaceholder("<value>"),
		},
	}

	suggestions := textInput.FlagSuggestions("", flags, nil)
	fmt.Printf("Text: %s, Description: %s\n", suggestions[0].Text, suggestions[0].Description)

	suggestions = textInput.FlagSuggestions("--", flags, nil)
	fmt.Printf("Text: %s, Description: %s\n", suggestions[0].Text, suggestions[0].Description)

	suggestions = textInput.FlagSuggestions("", flags,
		func(flagInput commandinput.FlagInput) commandinput.CommandMetadata[any] {
			return commandinput.CommandMetadata[any]{
				FlagArgPlaceholder:  flagInput.ArgPlaceholder,
				PreservePlaceholder: true,
			}
		})
	fmt.Printf(
		"Text: %s, Description: %s, Preserve Placeholder: %t\n",
		suggestions[0].Text,
		suggestions[0].Description,
		suggestions[0].Metadata.PreservePlaceholder,
	)

	// Output:
	// Text: -i, Description: refresh interval
	// Text: --interval, Description: refresh interval
	// Text: -i, Description: refresh interval, Preserve Placeholder: true
}

func ExampleModel_ParseUsage() {
	textInput := commandinput.New[commandinput.CommandMetadata[any]]()

	usage := `<mandatory arg> [optional arg] 'quoted arg' "double quoted arg" normal-arg`
	args, err := textInput.ParseUsage(usage)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n%s\n%s\n%s\n%s",
		args[0].Placeholder(),
		args[1].Placeholder(),
		args[2].Placeholder(),
		args[3].Placeholder(),
		args[4].Placeholder())

	// Output:
	// <mandatory arg>
	// [optional arg]
	// 'quoted arg'
	// "double quoted arg"
	// normal-arg
}
