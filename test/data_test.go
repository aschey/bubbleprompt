package test

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
)

var leftPadding = 2
var margin = 1

type cmdMetadata = commandinput.CmdMetadata

func suggestions(textInput *commandinput.Model[cmdMetadata]) []input.Suggestion[cmdMetadata] {
	return []input.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test desc", Metadata: commandinput.CmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[test placeholder1]", "[test placeholder2]"),
		}},
		{Text: "second-option", Description: "test desc2", Metadata: commandinput.CmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[test placeholder]"),
		}},
		{Text: "third-option", Description: "test desc3", Metadata: commandinput.CmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[flags]"),
		}},
		{Text: "fourth-option", Description: "test desc4"},
		{Text: "fifth-option", Description: "test desc5"},
		{Text: "sixth-option", Description: "test desc6"},
		{Text: "seventh-option", CompletionText: "completion text", Description: "test desc7"}}
}

func secondLevelSuggestions(textInput *commandinput.Model[cmdMetadata]) []input.Suggestion[cmdMetadata] {
	return []input.Suggestion[cmdMetadata]{
		{Text: "second-level", Description: "test desc", Metadata: commandinput.CmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[placeholder2]"),
			Level:          1,
		}},
	}
}

var flags = []commandinput.FlagInput{
	{Short: "t", Long: "test", Description: "test flag"},
}
