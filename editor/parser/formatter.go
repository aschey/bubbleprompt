package parser

import "github.com/aschey/bubbleprompt/editor"

type Formatter interface {
	Lex(input string, selectedToken *editor.Token) ([]FormatterToken, error)
}
