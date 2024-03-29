package lexer

type TokenType string

const (
	Undefined TokenType = ""

	// Single-character tokens.
	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"
	COMMA       = ","
	DOT         = "."
	MINUS       = "-"
	PLUS        = "+"
	SEMICOLON   = ";"
	SLASH       = "/"
	STAR        = "*"

	// One or two character tokens.
	BANG          = "!"
	BANG_EQUAL    = "!="
	EQUAL         = "="
	EQUAL_EQUAL   = "=="
	GREATER       = ">"
	GREATER_EQUAL = ">="
	LESS          = "<"
	LESS_EQUAL    = "<="

	// Literals.
	IDENTIFIER = "identifier"
	STRING     = "string"
	NUMBER     = "number"

	// Keywords.
	AND      = "and"
	CLASS    = "class"
	ELSE     = "else"
	FALSE    = "false"
	FUN      = "fun"
	FOR      = "for"
	IF       = "if"
	NIL      = "nil"
	OR       = "or"
	PRINT    = "print"
	RETURN   = "return"
	SUPER    = "super"
	THIS     = "this"
	TRUE     = "true"
	VAR      = "var"
	WHILE    = "while"
	BREAK    = "break"
	CONTINUE = "continue"
	EOF      = "EOF"
)

type Token struct {
	Type  TokenType
	Value interface{}
	Line  int
}

var Keywords = map[string]TokenType{
	"var":      VAR,
	"and":      AND,
	"or":       OR,
	"if":       IF,
	"else":     ELSE,
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
	"for":      FOR,
	"while":    WHILE,
	"fun":      FUN,
	"print":    PRINT,
	"return":   RETURN,
	"break":    BREAK,
	"continue": CONTINUE,
}
