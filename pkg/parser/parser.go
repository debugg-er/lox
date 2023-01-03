package parser

import (
	. "github.com/debugg-er/lox/pkg/common"
)

type Parser struct {
	current  int
	hadError bool
	tokens   []Token
}

func NewParser() *Parser {
	return &Parser{
		current:  0,
		hadError: false,
		tokens:   make([]Token, 0),
	}
}

func (p *Parser) Parse(tokens []Token) *Expr {
	p.tokens = tokens
	expr := p.expression()
	return expr
}

func (p *Parser) expression() *Expr {
	return p.binaryPrec(binRules, 0)
}

// `rules` parameter is an array of BinaryRUle that were defined
// with the priority go from highest to lowest accoding to its index
// `binaryPrec` should be passed zero value for `ruleIndex` parameter
func (p *Parser) binaryPrec(rules []BinaryRule, ruleIndex int) *Expr {
	if ruleIndex == len(rules) {
		return p.unary()
	}
	expr := p.binaryPrec(rules, ruleIndex+1)
	for {
		operator := p.match(rules[ruleIndex]...)
		if operator == nil {
			return expr
		}
		childPrec := p.binaryPrec(rules, ruleIndex+1)

		expr = &Expr{
			Type:     BINARY,
			Operator: operator,
			Left:     expr,
			Right:    childPrec,
		}
	}
}

func (p *Parser) unary() *Expr {
	operator := p.match(BANG, MINUS)
	if operator != nil {
		unaryExpr := p.unary()
		return &Expr{
			Type:     UNARY,
			Operator: operator,
			Left:     unaryExpr,
		}
	} else {
		return p.primary()
	}
}

func (p *Parser) primary() *Expr {
	token := p.advance()
	switch token.Type {
	case NUMBER, STRING, TRUE, FALSE, NIL:
		return NewLiteralExpr(token)
	case LEFT_PAREN:
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expected ')' after expression")
		return expr
	}
	// Give back the token
	p.current--
	return nil
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.current].Type == EOF
}

func (p *Parser) advance() *Token {
	token := p.tokens[p.current]
	p.current++
	return &token
}

func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) match(types ...TokenType) *Token {
	if p.isAtEnd() {
		return nil
	}
	for _, _type := range types {
		if p.peek().Type == _type {
			return p.advance()
		}
	}
	return nil
}

func (p *Parser) consume(tokenType TokenType, message string) {
	if p.peek().Type != tokenType {
		p.hadError = true
		Syncronize(p)
		return
	}
	p.advance()
}
