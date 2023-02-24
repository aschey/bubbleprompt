package commandinput

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/parser"
)

func buildCliParser(delimiterRegex string) *parser.ParticipleParser[statement] {
	lexer := lexer.MustStateful(lexer.Rules{
		"Root": {
			{Name: "Flag", Pattern: `\-{1,2}[^\s=\-]*`, Action: lexer.Push("Flag")},
			lexer.Include("Standard"),
			{Name: "Whitespace", Pattern: delimiterRegex},
		},
		"Standard": {
			{Name: "QuotedString", Pattern: `("[^"]*"?)|('[^']*'?)`},
			{Name: "String", Pattern: `[^\s]+`},
		},
		"Flag": {
			{Name: "Eq", Pattern: `\s*=\s*`, Action: lexer.Pop()},
			lexer.Include("Standard"),
			{Name: "FlagWhitespace", Pattern: delimiterRegex, Action: lexer.Pop()},
		},
	})
	participleParser := participle.MustBuild[statement](
		participle.Lexer(lexer),
		participle.Elide("Whitespace", "FlagWhitespace"))
	return parser.NewParticipleParser(participleParser)
}

type statement struct {
	Pos     lexer.Position
	Command ident `parser:"@@?"`
	Args    args  `parser:"@@"`
	Flags   flags `parser:"@@"`
	// Invalid input but this needs to be included to make the parser happy
	TrailingText []ident `parser:"@@?"`
}

type Statement struct {
	Command input.Token
	Args    []input.Token
	Flags   []Flag
}

func (s statement) toStatement() Statement {
	return Statement{
		Command: input.TokenFromPos(
			s.Command.Value,
			"command",
			0,
			s.Pos,
		),
		Args:  s.Args.toArgs(1),
		Flags: s.Flags.toFlags(len(s.Args.Value) + 1),
	}
}

type args struct {
	Pos   lexer.Position
	Value []ident `parser:"@@*"`
}

type Arg struct {
	Start int
	Value string
}

func (a args) toArgs(startIndex int) []input.Token {
	args := []input.Token{}
	for i, arg := range a.Value {
		args = append(
			args,
			input.TokenFromPos(arg.Value, "arg", startIndex+i, arg.Pos),
		)
	}
	return args
}

type flags struct {
	Pos   lexer.Position
	Value []flag `parser:"@@*"`
}

func (f flags) toFlags(startIndex int) []Flag {
	flags := []Flag{}
	for _, flag := range f.Value {
		flags = append(flags, flag.toFlag(startIndex))
		startIndex++
		if flag.Value != nil {
			startIndex++
		}
	}
	return flags
}

type flag struct {
	Pos   lexer.Position
	Name  string `parser:"( @Flag )"`
	Delim *delim `parser:"@@?"`
	Value *ident `parser:"@@?"`
}

func (f flag) toFlag(index int) Flag {
	name := input.TokenFromPos(f.Name, "flag", index, f.Pos)
	var value *input.Token = nil
	if f.Value != nil {
		v := input.TokenFromPos(
			f.Value.Value,
			"flagValue",
			index+1,
			f.Value.Pos,
		)
		value = &v
	}
	return Flag{
		Name:  name,
		Value: value,
	}
}

type Flag struct {
	Name  input.Token
	Value *input.Token
}

type delim struct {
	Pos   lexer.Position
	Value string `parser:"@Eq"`
}

type ident struct {
	Pos   lexer.Position
	Value string `parser:"( @QuotedString | @String )"`
}

func (i ident) ToToken(index int, tokenType string) input.Token {
	return input.TokenFromPos(i.Value, tokenType, index, i.Pos)
}
