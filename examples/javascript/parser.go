package main

import "github.com/alecthomas/participle/v2/lexer"

var lex = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Whitespace", Pattern: `\s+`},
	{Name: "String", Pattern: `"([^"]*"?)|('[^']*'?)`},
	{Name: "And", Pattern: `&&`},
	{Name: "Or", Pattern: `\|\|`},
	{Name: "Eq", Pattern: `===?`},
	{Name: "Number", Pattern: `[0-9]+(\.[0-9]*)*`},
	{Name: "Punct", Pattern: `[-\[!@#$%^&*()+_=\{\}\|:;"'<,>.?\/\]|]`},
	{Name: "Ident", Pattern: `[_a-zA-Z]+[_a-zA-Z0-9]*`},
})

type statement struct {
	Assignment *assignment `parser:"(@@"`
	Expression *expression `parser:"| @@)?"`
}

type assignment struct {
	Identifier string      `parser:" @Ident '=' "`
	Expression *expression `parser:"@@?"`
}

type indexer struct {
	OBracket   string      `parser:"@ '['"`
	Expression *expression `parser:"(@@"`
	CBracket   *string     `parser:"@ ']')?"`
}

type group struct {
	Expression *expression `parser:"'(' @@ ')'"`
}

type expression struct {
	Array        *array        `parser:"( @@"`
	Object       *object       `parser:"| @@"`
	Group        *group        `parser:"| @@"`
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
	Key   string      `parser:" @String | @Ident "`
	Delim *string     `parser:" @':'? "`
	Value *expression `parser:"  @@? "`
}

type object struct {
	Properties *[]keyValuePair `parser:" '{' (@@ (',' @@)*)* '}'? "`
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

type array struct {
	Values []expression `parser:" '[' (@@ ( ',' @@ )*)* ']' "`
}

type literal struct {
	Null    *string  `parser:" ( @( 'null' | 'undefined' ) "`
	Boolean *bool    `parser:" | @( 'true' | 'false' ) "`
	Str     *string  `parser:"| @String"`
	Number  *float64 `parser:"| @Number )"`
}
