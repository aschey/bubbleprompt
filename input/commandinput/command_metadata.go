package commandinput

import "github.com/aschey/bubbleprompt/suggestion"

// CommandMetadata defines the metadata that the [Model] uses to get information
// about the supplied [suggestion.Suggestion].
// You can extend this struct to provide additional metadata.
type CommandMetadata[T any] struct {
	// PositionalArgs is the list of positional args that this suggestion accepts.
	PositionalArgs []PositionalArg
	// ShowFlagPlaceholder is whether or not the input should display a placeholder
	// indicating that this command has flags available.
	ShowFlagPlaceholder bool
	// FlagArgPlaceholder is the placeholder
	FlagArgPlaceholder  FlagArgPlaceholder
	PreservePlaceholder bool
	Variadic            bool
	Children            []suggestion.Suggestion[CommandMetadata[T]]
	Extra               T
}

// MetadataFromPositionalArgs is a convenience function for creating a [CommandMetadata]
// from one or more [PositionalArg].
func MetadataFromPositionalArgs[T any](positionalArgs ...PositionalArg) CommandMetadata[T] {
	return CommandMetadata[T]{
		PositionalArgs: positionalArgs,
	}
}

func (c CommandMetadata[T]) GetChildren() []suggestion.Suggestion[CommandMetadata[T]] {
	return c.Children
}
