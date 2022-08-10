package parser

import "github.com/alecthomas/participle/v2"

type ParticipleParser[G any] struct {
	parser *participle.Parser[G]
}

func NewParticipleParser[G any](parser *participle.Parser[G]) *ParticipleParser[G] {
	return &ParticipleParser[G]{parser: parser}
}

func (p *ParticipleParser[G]) Lexer() Lexer {
	return NewParticipleLexer(p.parser.Lexer())
}
func (p *ParticipleParser[G]) Parse(input string) (*G, error) {
	return p.parser.ParseString("", input)
}
