package main

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/participle/v2"
	"github.com/aschey/bubbleprompt/parser/lexerbuilder"
)

var rules = []lexerbuilder.Rule{
	{Name: "Whitespace", Pattern: `\s+`, Type: chroma.Whitespace},
	{Name: "Grouping", Pattern: `[\(\)]`, Type: chroma.Punctuation},
	{Name: "String", Pattern: `"("[^"]*"?)|('[^']*'?)`, Type: chroma.String},
	{Name: "Number", Pattern: `\-?[0-9]+(\.[0-9]*)*`, Type: chroma.LiteralNumber},
	{Name: "Punct", Pattern: `[-\[!@#$%^&*+_=\{\}\|:;"'<,>.?\/\]|]`, Type: chroma.Punctuation},
	{Name: "Ident", Pattern: `[_a-zA-Z]+[_a-zA-Z0-9]*`, Type: chroma.Text},
}

var lex, styleLexer = lexerbuilder.NewLexerBuilder(rules).BuildLexers()

var participleParser = participle.MustBuild[statement](participle.Lexer(lex),
	participle.UseLookahead(20),
	participle.Elide("Whitespace", "Grouping"),
)

type statement struct {
	Assignment *assignment `parser:"(@@"`
	Expression *expression `parser:"| @@)?"`
}

type assignment struct {
	PropAccessor *propAccessor `parser:" ( @@ "`
	Identifier   *string       `parser:"| @Ident ) '=' "`
	Expression   *expression   `parser:"@@?"`
}

type indexer struct {
	OBracket   string      `parser:"@ '['"`
	Expression *expression `parser:"(@@"`
	CBracket   *string     `parser:"@ ']'?)?"`
}

type expression struct {
	Array        *array        `parser:"( @@"`
	Object       *object       `parser:"| @@"`
	PropAccessor *propAccessor `parser:"| @@"`
	Token        *token        `parser:"| @@)"`
	InfixOp      *infixOp      `parser:"(@@"`
	Expression   *expression   `parser:"@@)?"`
}

type token struct {
	Literal  *literal `parser:"@@"`
	Variable *string  `parser:"| @Ident"`
}

type keyValuePair struct {
	Key   string      `parser:" ( @String | @Ident ) "`
	Delim *string     `parser:" @':'? "`
	Value *expression `parser:"  @@? "`
}

type accessor struct {
	Indexer  *indexer  `parser:"(@@"`
	Delim    *string   `parser:"| (@ '.' "`
	Prop     *string   `parser:" @Ident?))"`
	Accessor *accessor `parser:"@@?"`
}

type propAccessor struct {
	Identifier string   `parser:" @Ident "`
	Accessor   accessor `parser:"@@"`
}

type infixOp struct {
	Op string `parser:" @( '+' | '-' | '*' | '/' | '||' | '&&' | '==' | '===' ) "`
}

type object struct {
	Properties []keyValuePair `parser:" '{' (@@ (',' @@)*)* '}'? "`
}

type array struct {
	Values []expression `parser:" '[' (@@ ( ',' @@? )*)* ']'? "`
}

type literal struct {
	Null    *string  `parser:" ( @( 'null' | 'undefined' ) "`
	Boolean *bool    `parser:" | @( 'true' | 'false' ) "`
	Str     *string  `parser:"| @String"`
	Number  *float64 `parser:"| @Number )"`
}
