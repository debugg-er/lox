package parser

import (
	. "github.com/debugg-er/lox/pkg/common"
)

type BinaryRule []TokenType

var binRules []BinaryRule = []BinaryRule{
	[]TokenType{EQUAL_EQUAL, BANG_EQUAL},                  // equality
	[]TokenType{GREATER, GREATER_EQUAL, LESS, LESS_EQUAL}, // comparison
	[]TokenType{PLUS, MINUS},                              // term
	[]TokenType{STAR, SLASH},                              // factor
}
