package commandinput_test

import (
	"fmt"

	"github.com/aschey/bubbleprompt/input/commandinput"
)

func ExampleModel_FlagSuggestions() {
	textInput := commandinput.New[commandinput.CmdMetadata]()
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
		func(flagInput commandinput.FlagInput) commandinput.CmdMetadata {
			return commandinput.CmdMetadata{
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
