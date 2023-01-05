package commandinput_test

import (
	"fmt"

	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
)

func ExampleMetadataFromPositionalArgs() {
	textInput := commandinput.New[commandinput.CommandMetadata]()
	commandMetadata := commandinput.MetadataFromPositionalArgs(textInput.NewPositionalArg("<arg1>"))

	suggestions := []suggestion.Suggestion[commandinput.CommandMetadata]{
		{Text: "test", Metadata: commandMetadata},
	}

	fmt.Println(suggestions[0].Metadata.GetPositionalArgs()[0].Placeholder())
	// Output: <arg1>
}
