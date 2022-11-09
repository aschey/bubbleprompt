package commandinput

type CmdMetadataAccessor interface {
	GetPositionalArgs() []PositionalArg
	GetFlagPlaceholder() FlagPlaceholder
	GetLevel() int
	GetPreservePlaceholder() bool
	GetShowFlagPlaceholder() bool
	Create(args []PositionalArg, placeholder FlagPlaceholder) CmdMetadataAccessor
}

type CmdMetadata struct {
	PositionalArgs      []PositionalArg
	ShowFlagPlaceholder bool
	FlagPlaceholder     FlagPlaceholder
	Level               int
	PreservePlaceholder bool
}

func MetadataFromPositionalArgs(positionalArgs ...PositionalArg) CmdMetadata {
	return CmdMetadata{
		PositionalArgs: positionalArgs,
	}
}

func (m CmdMetadata) Create(args []PositionalArg, placeholder FlagPlaceholder) CmdMetadataAccessor {
	return CmdMetadata{PositionalArgs: args, FlagPlaceholder: placeholder}
}

func (m CmdMetadata) GetPositionalArgs() []PositionalArg {
	return m.PositionalArgs
}

func (m CmdMetadata) GetFlagPlaceholder() FlagPlaceholder {
	return m.FlagPlaceholder
}

func (m CmdMetadata) GetShowFlagPlaceholder() bool {
	return m.ShowFlagPlaceholder
}

func (m CmdMetadata) GetLevel() int {
	return m.Level
}

func (m CmdMetadata) GetPreservePlaceholder() bool {
	return m.PreservePlaceholder
}
