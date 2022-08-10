package commandinput

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input/parser"
)

func (m *Model[T]) buildParser() {
	lexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "LongFlag", Pattern: `\-\-[^\s=\-]*`},
		{Name: "ShortFlag", Pattern: `\-[^\s=\-]*`},
		{Name: "Eq", Pattern: "="},
		{Name: "QuotedString", Pattern: `"[^"]*"`},
		{Name: `String`, Pattern: m.stringRegex.String()},
		{Name: "whitespace", Pattern: m.delimiterRegex.String()},
	})
	participleParser := participle.MustBuild[Statement](participle.Lexer(lexer))
	m.parser = parser.NewParticipleParser(participleParser)
}

type Statement struct {
	Pos     lexer.Position
	Command ident `parser:"@@?"`
	Args    args  `parser:"@@"`
	Flags   flags `parser:"@@"`
	// Invalid input but this needs to be included to make the parser happy
	TrailingText []ident `parser:"@@?"`
}

type args struct {
	Pos   lexer.Position
	Value []ident `parser:"@@*"`
}

type flags struct {
	Pos   lexer.Position
	Value []flag `parser:"@@*"`
}

type flag struct {
	Pos   lexer.Position
	Name  string `parser:"( @ShortFlag | @LongFlag )"`
	Delim *delim `parser:"@@?"`
	Value *ident `parser:"@@?"`
}

type delim struct {
	Pos   lexer.Position
	Value string `parser:"@Eq"`
}

type ident struct {
	Pos   lexer.Position
	Value string `parser:"( @QuotedString | @String )"`
}
