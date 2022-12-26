package parser

import (
	. "github.com/debugg-er/lox/pkg/common"
)

type Parser struct {
	tokens []Token
}

func (p *Parser) Parse() Expr {
	return Expr{}
}
