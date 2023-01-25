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

type RoundingBehavior int

const (
	RoundUp RoundingBehavior = iota
	RoundDown
)

func FindCurrentToken(
	runes []rune,
	tokens []Token,
	cursor int,
	roundingBehavior RoundingBehavior,
	isDelimiter func(s string, last Token) bool,
) Token {
	if len(tokens) > 0 {
		last := tokens[len(tokens)-1]
		index := len(tokens) - 1
		currentRuneIsDelimiter := cursor > 0 && len(runes) > 0 &&
			isDelimiter(string(runes[cursor-1]), last)
		if roundingBehavior == RoundUp && currentRuneIsDelimiter {
			// Haven't started a new token yet, but we have added a delimiter
			// so we'll consider the current token finished
			index++
		}
		// Check if cursor is at the end
		if cursor > last.End() {
			return Token{
				Start: cursor,
				Index: index,
			}
		}
	}
	for i := 0; i < len(tokens); i++ {
		if cursorInToken(tokens, cursor, i, roundingBehavior) {
			return tokens[i]
		}
	}
	return Token{Index: -1, Start: 0}
}

func cursorInToken(tokens []Token, cursor int, pos int, roundingBehavior RoundingBehavior) bool {
	isInToken := cursor >= tokens[pos].Start && cursor <= tokens[pos].End()
	if isInToken {
		return true
	}
	if roundingBehavior == RoundDown {
		if pos == len(tokens)-1 {
			return true
		}
		return cursor < tokens[pos+1].Start
	} else {
		if pos == 0 {
			return false
		}
		return cursor > tokens[pos-1].End() && cursor < tokens[pos].Start
	}
}
