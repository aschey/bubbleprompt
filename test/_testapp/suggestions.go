package testapp

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
)

type cmdMetadata = commandinput.CmdMetadata

var Suggestions []input.Suggestion[cmdMetadata] = []input.Suggestion[cmdMetadata]{
	{Text: "first-option", Description: "test desc", Metadata: commandinput.NewCmdMetadata([]commandinput.PositionalArg{commandinput.NewPositionalArg("[test placeholder]")}, commandinput.Placeholder{})},
	{Text: "second-option", Description: "test desc2"},
	{Text: "third-option", Description: "test desc3"},
	{Text: "fourth-option", Description: "test desc4"},
	{Text: "fifth-option", Description: "test desc5"},
	{Text: "sixth-option", Description: "test desc6"},
	{Text: "seventh-option", Description: "test desc7"}}
