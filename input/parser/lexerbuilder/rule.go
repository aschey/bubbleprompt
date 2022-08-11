package lexerbuilder

import "github.com/alecthomas/chroma/v2"

type Rule struct {
	Name    string
	Pattern string
	Type    chroma.Emitter
	Mutator chroma.Mutator
}
