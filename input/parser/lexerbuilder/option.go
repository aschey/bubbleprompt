package lexerbuilder

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Option func(builder *LexerBuilder) error

func WithLexerOptions(options ...lexer.Option) Option {
	return func(model *LexerBuilder) error {
		model.lexerOptions = options
		return nil
	}
}

func WithChromaConfig(config *chroma.Config) Option {
	return func(model *LexerBuilder) error {
		model.chromaConfig = config
		return nil
	}
}
