package parser

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type ParticipleLexer struct {
	definition lexer.Definition
}

func NewParticipleLexer(definition lexer.Definition) *ParticipleLexer {
	return &ParticipleLexer{definition: definition}
}

func (p *ParticipleLexer) Lex(input string) ([]Token, error) {
	lex, err := p.definition.Lex("", strings.NewReader(input))
	if err != nil {
		return nil, err
	}
	lexerTokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, err
	}
	if len(lexerTokens) > 0 {
		// Remove EOF token
		lexerTokens = lexerTokens[:len(lexerTokens)-1]
	}
	symbols := lexer.SymbolsByRune(p.definition)
	tokens := make([]Token, len(lexerTokens))
	for i, token := range lexerTokens {
		tokens[i] = Token{
			Start: token.Pos.Offset,
			Value: token.Value,
			Type:  symbols[token.Type],
			Index: i,
		}
	}

	return tokens, nil
}
