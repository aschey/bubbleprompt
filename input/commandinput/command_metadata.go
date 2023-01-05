package commandinput

// CommandMetadataAccessor defines the interface that the [Model] uses to get information about the supplied list of [suggestion.Suggestion].
type CommandMetadataAccessor interface {
	GetPositionalArgs() []PositionalArg
	GetFlagPlaceholder() FlagPlaceholder
	GetLevel() int
	GetPreservePlaceholder() bool
	GetShowFlagPlaceholder() bool
}

type CommandMetadata struct {
	PositionalArgs      []PositionalArg
	ShowFlagPlaceholder bool
	FlagPlaceholder     FlagPlaceholder
	Level               int
	PreservePlaceholder bool
}

func MetadataFromPositionalArgs(positionalArgs ...PositionalArg) CommandMetadata {
	return CommandMetadata{
		PositionalArgs: positionalArgs,
	}
}

func (m CommandMetadata) GetPositionalArgs() []PositionalArg {
	return m.PositionalArgs
}

func (m CommandMetadata) GetFlagPlaceholder() FlagPlaceholder {
	return m.FlagPlaceholder
}

func (m CommandMetadata) GetShowFlagPlaceholder() bool {
	return m.ShowFlagPlaceholder
}

func (m CommandMetadata) GetLevel() int {
	return m.Level
}

func (m CommandMetadata) GetPreservePlaceholder() bool {
	return m.PreservePlaceholder
}
