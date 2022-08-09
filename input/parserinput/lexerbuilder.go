package parserinput

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Rule struct {
	Name    string
	Pattern string
	Type    chroma.Emitter
	Mutator chroma.Mutator
}

type Rules []Rule

func (rules Rules) BuildLexers() ([]lexer.SimpleRule, []chroma.Rule) {
	lexerRules := make([]lexer.SimpleRule, len(rules))
	chromaRules := make([]chroma.Rule, len(rules))

	for i, rule := range rules {
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
