package parser

import "github.com/aschey/bubbleprompt/input"

type Lexer interface {
	Lex(input string) ([]input.Token, error)
}
