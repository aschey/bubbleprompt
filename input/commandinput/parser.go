package commandinput

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/parser"
)

func (m *Model[T]) buildParser() {
	lexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "LongFlag", Pattern: `\-\-[^\s=\-]*`},
		{Name: "ShortFlag", Pattern: `\-[^\s=\-]*`},
		{Name: "Eq", Pattern: "="},
		{Name: "QuotedString", Pattern: `("[^"]*"?)|('[^']*'?)`},
		{Name: `String`, Pattern: `[^\-\s][^\s]*`},
		{Name: "whitespace", Pattern: m.delimiterRegex.String()},
	})
	participleParser := participle.MustBuild[statement](participle.Lexer(lexer))
	m.parser = parser.NewParticipleParser(participleParser)
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
		Command: input.TokenFromPos(s.Command.Value, "command", 0, s.Pos), //TokenValue{value: s.Command.Value},
		Args:    s.Args.toArgs(1),
		Flags:   s.Flags.toFlags(len(s.Args.Value) + 1),
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
		args = append(args, input.TokenFromPos(arg.Value, "arg", startIndex+i, arg.Pos)) //TokenValue{value: arg.Value})
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
	Name  string `parser:"( @ShortFlag | @LongFlag )"`
	Delim *delim `parser:"@@?"`
	Value *ident `parser:"@@?"`
}

func (f flag) toFlag(index int) Flag {
	name := input.TokenFromPos(f.Name, "flag", index, f.Pos)
	var value *input.Token = nil
	if f.Value != nil {
		v := input.TokenFromPos(f.Value.Value, "flagValue", index+1, f.Value.Pos) // {value: f.Value.Value}
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
