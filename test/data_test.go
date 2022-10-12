package test

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
)

var leftPadding = 2
var margin = 1

type cmdMetadata = commandinput.CmdMetadata

var suggestions []input.Suggestion[cmdMetadata] = []input.Suggestion[cmdMetadata]{
	{Text: "first-option", Description: "test desc", Metadata: commandinput.CmdMetadata{
		PositionalArgs: []commandinput.PositionalArg{
			commandinput.NewPositionalArg("[test placeholder1]"),
			commandinput.NewPositionalArg("[test placeholder2]"),
		},
	}},
	{Text: "second-option", Description: "test desc2", Metadata: commandinput.CmdMetadata{
		PositionalArgs: []commandinput.PositionalArg{
			commandinput.NewPositionalArg("[test placeholder]"),
		},
	}},
	{Text: "third-option", Description: "test desc3"},
	{Text: "fourth-option", Description: "test desc4"},
	{Text: "fifth-option", Description: "test desc5"},
	{Text: "sixth-option", Description: "test desc6"},
	{Text: "seventh-option", CompletionText: "completion text", Description: "test desc7"}}

var secondLevelSuggestions []input.Suggestion[cmdMetadata] = []input.Suggestion[cmdMetadata]{
	{Text: "second-level", Description: "test desc", Metadata: commandinput.CmdMetadata{
		PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("[placeholder2]")},
		Level:          1,
	}},
}
