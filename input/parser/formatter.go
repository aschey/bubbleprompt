package parser

type Formatter interface {
	Lex(input string) ([]FormatterToken, error)
}
