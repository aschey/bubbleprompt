package commandinput

type CmdMetadataAccessor interface {
	GetPositionalArgs() []PositionalArg
	GetFlagPlaceholder() Placeholder
	GetLevel() int
	GetPreservePlaceholder() bool
	GetShowFlagPlaceholder() bool
	Create(args []PositionalArg, placeholder Placeholder) CmdMetadataAccessor
}

type CmdMetadata struct {
	PositionalArgs      []PositionalArg
	ShowFlagPlaceholder bool
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

func (m CmdMetadata) GetShowFlagPlaceholder() bool {
	return m.ShowFlagPlaceholder
}

func (m CmdMetadata) GetLevel() int {
	return m.Level
}

func (m CmdMetadata) GetPreservePlaceholder() bool {
	return m.PreservePlaceholder
}
