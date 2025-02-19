package lexerbuilder

import (
	"github.com/alecthomas/chroma/v2"
)

type Option func(builder *LexerBuilder)

func WithChromaConfig(config *chroma.Config) Option {
	return func(model *LexerBuilder) {
		model.chromaConfig = config
	}
}
