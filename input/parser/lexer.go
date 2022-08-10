package parser

type Lexer interface {
	Lex(input string) ([]Token, error)
}
