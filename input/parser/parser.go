package parser

type Parser[G any] interface {
	Lexer() Lexer
	Parse(input string) (*G, error)
}
