package lexerbuilder

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Option func(builder *LexerBuilder)

func WithLexerOptions(options ...lexer.Option) Option {
	return func(model *LexerBuilder) {
		model.lexerOptions = options
	}
}

func WithChromaConfig(config *chroma.Config) Option {
	return func(model *LexerBuilder) {
		model.chromaConfig = config
	}
}
