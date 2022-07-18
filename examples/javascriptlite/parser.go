package main

type statement struct {
	Assignment *assignment `parser:"(@@"`
	Expression *expression `parser:"| @@)?"`
}

type assignment struct {
	Identifier identifier `parser:" @@ '=' "`
	Expression expression `parser:"@@?"`
}

type identifier struct {
	Variable string      `parser:"@Ident"`
	Accessor *expression `parser:" ('[' @@ ']?')? "`
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
	Literal  *literal    `parser:"@@"`
	Variable *identifier `parser:"| @@"`
}

type keyValuePair struct {
	Key   string     `parser:" @String "`
	Value expression `parser:" ':' @@ "`
}

type object struct {
	Properties *[]keyValuePair `parser:" '{' (@@ (',' @@)*)* '}' "`
}

type propAccessor struct {
	Identifier identifier    `parser:" @@ '.' "`
	Accessor   *propAccessor `parser:"( @@"`
	Prop       *string       `parser:" | @Ident )?"`
}

type infixOp struct {
	Op string `parser:" '+' | '-' | '*' | '/' | '||' | '&&' | '==' | '===' "`
}

type array struct {
	Values []expression `parser:" '[' (@@ ( ',' @@ )*)* ']' "`
}

type literal struct {
	Null    *string  `parser:" ( ( 'null' | 'undefined' ) "`
	Boolean *bool    `parser:" | ( 'true' | 'false' ) "`
	Str     *string  `parser:"| @String"`
	Number  *float64 `parser:"| @Number )"`
}
