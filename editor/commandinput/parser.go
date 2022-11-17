package commandinput

import (
	"encoding/json"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/editor/parser"
)

func (m *Model[T]) buildParser() {
	lexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "LongFlag", Pattern: `\-\-[^\s=\-]*`},
		{Name: "ShortFlag", Pattern: `\-[^\s=\-]*`},
		{Name: "Eq", Pattern: "="},
		{Name: "QuotedString", Pattern: `"[^"]*"`},
		{Name: `String`, Pattern: `[^\-\s][^\s]*`},
		{Name: "whitespace", Pattern: m.delimiterRegex.String()},
	})
	participleParser := participle.MustBuild[statement](participle.Lexer(lexer))
	m.parser = parser.NewParticipleParser(participleParser)
}

type TokenValue struct {
	value string
}

func (t TokenValue) RawValue() string {
	return t.value
}

func (t TokenValue) Value() string {
	if strings.HasPrefix(t.value, "\"") {
		var dest string
		err := json.Unmarshal([]byte(t.value), &dest)
		if err != nil {
			return t.value
		}
		return dest
	}
	return strings.ReplaceAll(t.value, `\"`, `"`)
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
	Start   int
	Command TokenValue
	Args    []TokenValue
	Flags   []Flag
}

func (s statement) toStatement() Statement {
	return Statement{
		Start:   s.Pos.Column - 1,
		Command: TokenValue{value: s.Command.Value},
		Args:    s.Args.toArgs(),
		Flags:   s.Flags.toFlags(),
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

func (a args) toArgs() []TokenValue {
	args := []TokenValue{}
	for _, arg := range a.Value {
		args = append(args, TokenValue{value: arg.Value})
	}
	return args
}

type flags struct {
	Pos   lexer.Position
	Value []flag `parser:"@@*"`
}

func (f flags) toFlags() []Flag {
	flags := []Flag{}
	for _, flag := range f.Value {
		flags = append(flags, flag.toFlag())
	}
	return flags
}

type flag struct {
	Pos   lexer.Position
	Name  string `parser:"( @ShortFlag | @LongFlag )"`
	Delim *delim `parser:"@@?"`
	Value *ident `parser:"@@?"`
}

func (f flag) toFlag() Flag {
	var value *TokenValue = nil
	if f.Value != nil {
		value = &TokenValue{value: f.Value.Value}
	}
	return Flag{
		Start: f.Pos.Column - 1,
		Name:  f.Name,
		Value: value,
	}
}

type Flag struct {
	Start int
	Name  string
	Value *TokenValue
}

type delim struct {
	Pos   lexer.Position
	Value string `parser:"@Eq"`
}

type ident struct {
	Pos   lexer.Position
	Value string `parser:"( @QuotedString | @String )"`
}

func (i ident) ToToken(index int, tokenType string) editor.Token {
	return editor.TokenFromPos(i.Value, tokenType, index, i.Pos)
}
