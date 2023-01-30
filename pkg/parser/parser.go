package parser

import (
	//lint:ignore ST1001 that's what we want
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

func (p *Parser) Parse(tokens []Token) ([]Stmt, []*Error) {
	p.tokens = tokens
	statements := make([]Stmt, 0)
	errors := make([]*Error, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errors = append(errors, err)
			p.synchronize()
			continue
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return statements, errors
}

func (p *Parser) declaration() (Stmt, *Error) {
	if p.match(VAR) != nil {
		return p.varDecl()
	}
	return p.statement()
}

func (p *Parser) varDecl() (Stmt, *Error) {
	token := p.advance()
	if token.Type != IDENTIFIER {
		return nil, NewError(token, "Expected variable name.")
	}
	var initilizer *Expr = nil
	if p.match(EQUAL) != nil {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		initilizer = expr
	}
	if err := p.consume(SEMICOLON, "Expected ';' after expression"); err != nil {
		return nil, err
	}
	return &VarStmt{token, initilizer}, nil
}

func (p *Parser) statement() (Stmt, *Error) {
	if p.match(PRINT) != nil {
		return p.printStmt()
	}
	if p.match(LEFT_BRACE) != nil {
		return p.blockStmt()
	}
	if p.match(IF) != nil {
		return p.ifStmt()
	}
	return p.exprStmt()
}

func (p *Parser) ifStmt() (Stmt, *Error) {
	if err := p.consume(LEFT_PAREN, "Expected '(' after if"); err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(RIGHT_PAREN, "Expected ')' after expression"); err != nil {
		return nil, err
	}
	trueStmt, err := p.statement()
	if err != nil {
		return nil, err
	}
	var falseStmt Stmt = nil
	if p.match(ELSE) != nil {
		falseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &IfStmt{
		condition: expr,
		trueStmt:  trueStmt,
		falseStmt: falseStmt,
	}, nil
}

func (p *Parser) blockStmt() (Stmt, *Error) {
	declarations := make([]Stmt, 0)
	for !p.isAtEnd() && p.peek().Type != RIGHT_BRACE {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		declarations = append(declarations, declaration)
	}
	if err := p.consume(RIGHT_BRACE, "Expected '}' after block"); err != nil {
		return nil, err
	}
	return &BlockStmt{
		declarations: declarations,
	}, nil
}

func (p *Parser) printStmt() (Stmt, *Error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err = p.consume(SEMICOLON, "Expected ';' after value"); err != nil {
		return nil, err
	}
	return &PrintStmt{expr}, nil
}

func (p *Parser) exprStmt() (Stmt, *Error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err = p.consume(SEMICOLON, "Expected ';' after expression"); err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, nil
	}
	return &ExprStmt{expr}, nil
}

func (p *Parser) expression() (*Expr, *Error) {
	return p.assignment()
}

func (p *Parser) assignment() (*Expr, *Error) {
	expr, err := p.binaryPrec(binRules, LOGICAL_OR)
	if err != nil {
		return nil, err
	}

	if equal := p.match(EQUAL); equal != nil {
		if expr.Type != VARIABLE {
			return nil, NewError(equal, "Invalid assignment target.")
		}
		assignment, err := p.assignment()
		if err != nil {
			return nil, err
		}
		return &Expr{
			Type:        ASSIGN,
			Var:         expr.Var,
			AssignValue: assignment,
		}, nil
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
		return NewPrimaryExpr(token), nil
	case LEFT_PAREN:
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if err = p.consume(RIGHT_PAREN, "Expected ')' after expression"); err != nil {
			return nil, err
		}
		return expr, nil
	case IDENTIFIER:
		return &Expr{
			Type: VARIABLE,
			Var:  token,
		}, nil
	default:
		// Give back the token if don't match any precedence
		p.current--
		return nil, nil
	}
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

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
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

func (p *Parser) synchronize() {
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
