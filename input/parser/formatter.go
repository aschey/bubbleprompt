package parser

type Formatter interface {
	Lex(input string, selectedToken *Token) ([]FormatterToken, error)
}
