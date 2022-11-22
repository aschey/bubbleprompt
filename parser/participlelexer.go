package parser

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
)

type ParticipleLexer struct {
	definition lexer.Definition
}

func NewParticipleLexer(definition lexer.Definition) *ParticipleLexer {
	return &ParticipleLexer{definition: definition}
}

func (p *ParticipleLexer) Lex(inputStr string) ([]input.Token, error) {
	lex, err := p.definition.Lex("", strings.NewReader(inputStr))
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
	tokens := make([]input.Token, len(lexerTokens))
	for i, token := range lexerTokens {
		tokens[i] = input.TokenFromPos(token.Value, symbols[token.Type], i, token.Pos)
	}

	return tokens, nil
}
