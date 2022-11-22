package parser

import "github.com/aschey/bubbleprompt/input"

type Formatter interface {
	Lex(input string, selectedToken *input.Token) ([]FormatterToken, error)
}
