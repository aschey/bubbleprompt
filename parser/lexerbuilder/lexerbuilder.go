package lexerbuilder

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type LexerBuilder struct {
	rules        []Rule
	lexerOptions []lexer.Option
	chromaConfig *chroma.Config
}

func NewLexerBuilder(rules []Rule, options ...Option) *LexerBuilder {
	builder := &LexerBuilder{rules: rules}
	for _, option := range options {
		option(builder)
	}
	return builder
}

func (b *LexerBuilder) BuildRuleLists() ([]lexer.SimpleRule, []chroma.Rule) {
	lexerRules := make([]lexer.SimpleRule, len(b.rules))
	chromaRules := make([]chroma.Rule, len(b.rules))

	for i, rule := range b.rules {
		lexerRules[i] = lexer.SimpleRule{
			Name:    rule.Name,
			Pattern: rule.Pattern,
		}

		chromaRules[i] = chroma.Rule{
			Pattern: rule.Pattern,
			Type:    rule.Type,
			Mutator: rule.Mutator,
		}
	}

	return lexerRules, chromaRules
}

func (b *LexerBuilder) BuildLexers() (*lexer.StatefulDefinition, *chroma.RegexLexer) {
	lexerRules, chromaRules := b.BuildRuleLists()
	lex := lexer.MustSimple(lexerRules, b.lexerOptions...)

	styleLexer := chroma.MustNewLexer(b.chromaConfig,
		func() chroma.Rules { return chroma.Rules{"root": chromaRules} })
	return lex, styleLexer
}
