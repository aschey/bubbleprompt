package commandinput

// CommandMetadataAccessor defines the interface that the [Model] uses
// to get information about the supplied [suggestion.Suggestion].
type CommandMetadataAccessor interface {
	// GetPositionalArgs returns the list of PositionalArgs for the associated suggestion.
	GetPositionalArgs() []PositionalArg
	// GetFlagArgPlaceholder returns the placeholder for the flag's argument.
	GetFlagArgPlaceholder() FlagArgPlaceholder
	// TODO: Figure out why we need this.
	GetPreservePlaceholder() bool
	// GetShowFlagPlaceholder returns whether or not to show the placeholder
	// indicating this suggestion has available flags.
	GetShowFlagPlaceholder() bool
}

// CommandMetadata defines the metadata that the [Model] uses to get information
// about the supplied [suggestion.Suggestion].
// You can extend this struct to provide additional metadata.
type CommandMetadata struct {
	// PositionalArgs is the list of positional args that this suggestion accepts.
	PositionalArgs []PositionalArg
	// ShowFlagPlaceholder is whether or not the input should display a placeholder
	// indicating that this command has flags available.
	ShowFlagPlaceholder bool
	// FlagArgPlaceholder is the placeholder
	FlagArgPlaceholder  FlagArgPlaceholder
	PreservePlaceholder bool
}

// MetadataFromPositionalArgs is a convenience function for creating a [CommandMetadata]
// from one or more [PositionalArg].
func MetadataFromPositionalArgs(positionalArgs ...PositionalArg) CommandMetadata {
	return CommandMetadata{
		PositionalArgs: positionalArgs,
	}
}

// GetPositionalArgs returns the list of [PositionalArg] for the associated [suggestion.Suggestion].
func (m CommandMetadata) GetPositionalArgs() []PositionalArg {
	return m.PositionalArgs
}

// GetFlagArgPlaceholder returns the [FlagPlaceholder] for the flag's argument.
func (m CommandMetadata) GetFlagArgPlaceholder() FlagArgPlaceholder {
	return m.FlagArgPlaceholder
}

// GetShowFlagPlaceholder returns whether or not to show the placeholder
// indicating this [suggestion.Suggestion] has available flags.
func (m CommandMetadata) GetShowFlagPlaceholder() bool {
	return m.ShowFlagPlaceholder
}

// TODO: Figure out why we need this.
func (m CommandMetadata) GetPreservePlaceholder() bool {
	return m.PreservePlaceholder
}
