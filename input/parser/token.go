package parser

type Token struct {
	Start int
	Type  string
	Value string
	Index int
}

func (t Token) End() int {
	return t.Start + len(t.Value)
}
