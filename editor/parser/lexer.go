package parser

import "github.com/aschey/bubbleprompt/editor"

type Lexer interface {
	Lex(input string) ([]editor.Token, error)
}
