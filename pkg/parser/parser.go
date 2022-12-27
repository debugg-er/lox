package parser

import (
	"fmt"

	. "github.com/debugg-er/lox/pkg/common"
)

type Parser struct {
	current int
	tokens  []Token
}

func NewParser() *Parser {
	return &Parser{
		current: 0,
		tokens:  make([]Token, 0),
	}
}

func (p *Parser) Parse(tokens []Token) *Expr {
	p.tokens = tokens
	return p.expression()
}

func (p *Parser) expression() *Expr {
	return p.equality()
}

func (p *Parser) equality() *Expr {
	return p.precedence(p.comparison, BANG_EQUAL, EQUAL_EQUAL)
}

func (p *Parser) comparison() *Expr {
	return p.precedence(p.term, GREATER, GREATER_EQUAL, LESS, LESS_EQUAL)
}

func (p *Parser) term() *Expr {
	return p.precedence(p.factor, PLUS, MINUS)
}

func (p *Parser) factor() *Expr {
	return p.precedence(p.unary, STAR, SLASH)
}

func (p *Parser) precedence(childPrec func() *Expr, types ...TokenType) *Expr {
	expr := childPrec()
	for {
		operator := p.matches(types...)
		if operator == Undefined {
			return expr
		}

		expr = &Expr{
			Type:     BINARY,
			Operator: operator,
			Left:     expr,
			Right:    childPrec(),
		}
	}
}

func (p *Parser) unary() *Expr {
	operator := p.matches(BANG, MINUS)
	fmt.Println(operator)
	if operator != Undefined {
		return &Expr{
			Type:     UNARY,
			Operator: operator,
			Left:     p.unary(),
		}
	} else {
		return p.primary()
	}
}

func (p *Parser) primary() *Expr {
	token := p.advance()
	switch token.Type {
	case NUMBER:
		return &Expr{
			Type:    LITERAL,
			Literal: token.Value,
		}
	case STRING:
		return &Expr{
			Type:    LITERAL,
			Literal: token.Value,
		}
	case TRUE:
		return &Expr{
			Type:    LITERAL,
			Literal: true,
		}
	case FALSE:
		return &Expr{
			Type:    LITERAL,
			Literal: false,
		}
	case NIL:
		return &Expr{
			Type:    LITERAL,
			Literal: nil,
		}
	case LEFT_PAREN:
		return &Expr{
			Type: GROUPING,
			Left: p.expression(),
		}
		// handle missing ')' error
	}
	// Give back the token
	p.current--
	return nil
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.current].Type == EOF
}

func (p *Parser) advance() Token {
	token := p.tokens[p.current]
	p.current++
	return token
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) matches(types ...TokenType) TokenType {
	if p.isAtEnd() {
		return Undefined
	}
	for _, _type := range types {
		if p.peek().Type == _type {
			p.advance()
			return _type
		}
	}
	return Undefined
}
