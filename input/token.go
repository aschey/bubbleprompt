package input

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Token struct {
	Start int
	Type  string
	Value string
	Index int
}

func (t Token) End() int {
	return t.Start + len([]rune(t.Value))
}

func TokenFromPos(value string, tokenType string, index int, pos lexer.Position) Token {
	start := pos.Column
	// Need to subtract 1 from the column to get the 0-based offset of the token
	if start > 0 {
		start--
	}
	return Token{
		Value: value,
		Type:  tokenType,
		Start: start,
	}
}
