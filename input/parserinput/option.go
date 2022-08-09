package parserinput

import "github.com/alecthomas/chroma/v2"

type Option func(model *LexerModel) error

func WithDelimiterTokens(tokens ...string) Option {
	return func(model *LexerModel) error {
		model.delimiterTokens = tokens
		return nil
	}
}

func WithDelimiters(delimiters ...string) Option {
	return func(model *LexerModel) error {
		model.delimiters = delimiters
		return nil
	}
}

func WithStyle(styleLexer chroma.Lexer, style chroma.Style) Option {
	return func(model *LexerModel) error {
		model.styleLexer = styleLexer
		model.style = &style
		return nil
	}
}
