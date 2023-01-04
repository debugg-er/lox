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

func (p *Parser) Parse(tokens []Token) []Stmt {
	p.tokens = tokens
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			p.Errors = append(p.Errors, *err)
			p.Synchronize()
		} else {
			statements = append(statements, *stmt)
		}
	}
	return statements
}

func (p *Parser) statement() (*Stmt, *Error) {
	var stmtType StmtType
	if p.peek().Type == PRINT {
		stmtType = PRINT_STMM
		p.advance()
	} else {
		stmtType = EXPR_STMT
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if err = p.consume(SEMICOLON, "Expected ;"); err != nil {
		return nil, err
	}
	return &Stmt{
		Type: stmtType,
		Expr: expr,
	}, nil
}

func (p *Parser) expression() (*Expr, *Error) {
	expr, err := p.binaryPrec(binRules, 0)
	if err != nil {
		return nil, err
	}
	if _, err := expr.Evaluate(); err != nil {
		return nil, err
	}
	return expr, nil
}

// `rules` parameter is an array of BinaryRule that were defined
// with the priority go from highest to lowest accoding to its index.
// `binaryPrec` should be passed zero value for `ruleIndex` parameter
func (p *Parser) binaryPrec(rules []BinaryRule, ruleIndex int) (*Expr, *Error) {
	if ruleIndex == len(rules) {
		return p.unary()
	}
	expr, err := p.binaryPrec(rules, ruleIndex+1)
	if err != nil {
		return nil, err
	}
	for {
		operator := p.match(rules[ruleIndex]...)
		if operator == nil {
			return expr, nil
		}
		childPrec, err := p.binaryPrec(rules, ruleIndex+1)
		if err != nil {
			return nil, err
		}

		expr = &Expr{
			Type:     BINARY,
			Operator: operator,
			Left:     expr,
			Right:    childPrec,
		}
	}
}

func (p *Parser) unary() (*Expr, *Error) {
	operator := p.match(BANG, MINUS)
	if operator == nil {
		return p.primary()
	}
	unaryExpr, err := p.unary()
	if err != nil {
		return nil, err
	}
	return &Expr{
		Type:     UNARY,
		Operator: operator,
		Left:     unaryExpr,
	}, nil
}

func (p *Parser) primary() (*Expr, *Error) {
	token := p.advance()
	switch token.Type {
	case NUMBER, STRING, TRUE, FALSE, NIL:
		return NewLiteralExpr(token), nil
	case LEFT_PAREN:
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if err = p.consume(RIGHT_PAREN, "Expected ')' after expression"); err != nil {
			return nil, err
		}
		return expr, nil
	}
	// Give back the token
	p.current--
	return nil, nil
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
	previous := p.peek()

	for !p.isAtEnd() {
		if previous.Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		previous = p.advance()
	}
}

func (p *Parser) consume(tokenType TokenType, message string) *Error {
	if p.isAtEnd() {
		return nil
	}
	if p.peek().Type != tokenType {
		return NewError(p.peek(), message)
	}
	p.advance()
	return nil
}
