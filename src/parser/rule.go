package parser

import (
	l "github.com/debugg-er/lox/src/lexer"
)

type BinaryRule []l.TokenType

const (
	LOGICAL_OR int = iota
	LOGICAL_AND
	EQUALITY
	COMPARISON
	TERM
	FACTOR
)

var binRules []BinaryRule = []BinaryRule{
	[]l.TokenType{l.OR},                        // logical_or
	[]l.TokenType{l.AND},                       // logical_and
	[]l.TokenType{l.EQUAL_EQUAL, l.BANG_EQUAL}, // equality
	[]l.TokenType{l.GREATER, l.GREATER_EQUAL, l.LESS, l.LESS_EQUAL}, // comparison
	[]l.TokenType{l.PLUS, l.MINUS},                                  // term
	[]l.TokenType{l.STAR, l.SLASH},                                  // factor
}
