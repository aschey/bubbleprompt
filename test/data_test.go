package test

import (
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
)

var (
	leftPadding = 2
	margin      = 1
)

type cmdMetadata = commandinput.CommandMetadata[any]

func suggestions(textInput *commandinput.Model[cmdMetadata]) []suggestion.Suggestion[cmdMetadata] {
	return []suggestion.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test desc", Metadata: cmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs(
				"[test placeholder1]",
				"[test placeholder2]",
			),
		}},
		{Text: "second-option", Description: "test desc2", Metadata: cmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[test placeholder]"),
		}},
		{Text: "third-option", Description: "test desc3", Metadata: cmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[flags]"),
		}},
		{Text: "fourth-option", Description: "test desc4"},
		{Text: "fifth-option", Description: "test desc5"},
		{Text: "sixth-option", Description: "test desc6"},
		{Text: "seventh-option", SuggestionText: "suggestion text", Description: "test desc7"},
	}
}

func secondLevelSuggestions(
	textInput *commandinput.Model[cmdMetadata],
) []suggestion.Suggestion[cmdMetadata] {
	return []suggestion.Suggestion[cmdMetadata]{
		{Text: "second-level", Description: "test desc", Metadata: cmdMetadata{
			PositionalArgs: textInput.NewPositionalArgs("[placeholder2]"),
		}},
	}
}

var flags = []commandinput.FlagInput{
	{Short: "t", Long: "test", Description: "test flag"},
}
