package parser

import (
	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type BinaryRule []TokenType

const (
	LOGICAL_OR int = iota
	LOGICAL_AND
	EQUALITY
	COMPARISON
	TERM
	FACTOR
)

var binRules []BinaryRule = []BinaryRule{
	[]TokenType{OR},                      // logical_or
	[]TokenType{AND},                     // logical_and
	[]TokenType{EQUAL_EQUAL, BANG_EQUAL}, // equality
	[]TokenType{GREATER, GREATER_EQUAL, LESS, LESS_EQUAL}, // comparison
	[]TokenType{PLUS, MINUS},                              // term
	[]TokenType{STAR, SLASH},                              // factor
}
