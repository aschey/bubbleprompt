package commandinput

type CmdMetadataAccessor interface {
	GetPositionalArgs() []PositionalArg
	GetFlagPlaceholder() Placeholder
	GetLevel() int
	GetPreservePlaceholder() bool
	GetHasFlags() bool
	Create(args []PositionalArg, placeholder Placeholder) CmdMetadataAccessor
}

type CmdMetadata struct {
	PositionalArgs      []PositionalArg
	HasFlags            bool
	FlagPlaceholder     Placeholder
	Level               int
	PreservePlaceholder bool
}

func (m CmdMetadata) Create(args []PositionalArg, placeholder Placeholder) CmdMetadataAccessor {
	return CmdMetadata{PositionalArgs: args, FlagPlaceholder: placeholder}
}

func (m CmdMetadata) GetPositionalArgs() []PositionalArg {
	return m.PositionalArgs
}

func (m CmdMetadata) GetFlagPlaceholder() Placeholder {
	return m.FlagPlaceholder
}

func (m CmdMetadata) GetHasFlags() bool {
	return m.HasFlags
}

func (m CmdMetadata) GetLevel() int {
	return m.Level
}

func (m CmdMetadata) GetPreservePlaceholder() bool {
	return m.PreservePlaceholder
}
