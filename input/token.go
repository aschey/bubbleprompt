package input

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/internal"
)

type Token struct {
	Start int
	Type  string
	Value string
	Index int
}

func (t Token) Unquote() string {
	if strings.HasPrefix(t.Value, `"`) || strings.HasPrefix(t.Value, `\"`) {
		return internal.Unescape(t.Value, `"`)
	}
	if strings.HasPrefix(t.Value, `'`) || strings.HasPrefix(t.Value, `\'`) {
		return internal.Unescape(t.Value, `'`)
	}
	return t.Value
}

func (t Token) Unescape(wrapper string) string {
	return internal.Unescape(t.Value, wrapper)
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
