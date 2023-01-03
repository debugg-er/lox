package parser

import (
	. "github.com/debugg-er/lox/pkg/common"
)

type Parser struct {
	Errors  []Error
	current int
	tokens  []Token
}

func NewParser() *Parser {
	return &Parser{
		current: 0,
		Errors:  make([]Error, 0),
		tokens:  make([]Token, 0),
	}
}

func (p *Parser) Parse(tokens []Token) *Expr {
	p.tokens = tokens
	expr := p.expression()
	if _, err := expr.Evaluate(); err != nil {
		p.Errors = append(p.Errors, *err)
	}
	return expr
}

func (p *Parser) expression() *Expr {
	return p.binaryPrec(binRules, 0)
}

// `rules` parameter is an array of BinaryRule that were defined
// with the priority go from highest to lowest accoding to its index.
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
	if operator == nil {
		return p.primary()
	}
	unaryExpr := p.unary()
	return &Expr{
		Type:     UNARY,
		Operator: operator,
		Left:     unaryExpr,
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
	if p.isAtEnd() {
		return nil
	}
	token := p.tokens[p.current]
	p.current++
	return &token
}

func (p *Parser) peek() *Token {
	if p.isAtEnd() {
		return nil
	}
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

func (p *Parser) Synchronize() {
	current := p.advance()

	for !p.isAtEnd() {
		if current.Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		current = p.advance()
	}
}

func (p *Parser) consume(tokenType TokenType, message string) {
	if p.peek().Type != tokenType {
		p.Errors = append(p.Errors, *NewError(p.peek(), message))
		p.Synchronize()
		return
	}
	p.advance()
}
