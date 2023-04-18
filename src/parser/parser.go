package parser

import (
	l "github.com/debugg-er/lox/src/lexer"
)

type Parser struct {
	current int
	tokens  []l.Token
}

func NewParser() *Parser {
	return &Parser{
		current: 0,
		tokens:  make([]l.Token, 0),
	}
}

func (p *Parser) Parse(tokens []l.Token) ([]Stmt, []error) {
	p.tokens = tokens
	statements := make([]Stmt, 0)
	errors := make([]error, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errors = append(errors, err)
			p.synchronize()
			continue
		}
		if stmt == nil {
			continue
		}
		if branchingErrors := verifyBranching(stmt); branchingErrors != nil {
			errors = append(errors, branchingErrors...)
			continue
		}
		statements = append(statements, stmt)
	}
	return statements, errors
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(l.VAR) != nil {
		return p.varDecl()
	}
	return p.statement()
}

func (p *Parser) varDecl() (Stmt, error) {
	token := p.advance()
	if token.Type != l.IDENTIFIER {
		return nil, NewParserError(token, "Expected variable name.")
	}
	var initilizer Expr = nil
	if p.match(l.EQUAL) != nil {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		initilizer = expr
	}
	if err := p.consume(l.SEMICOLON, "Expected ';' after expression"); err != nil {
		return nil, err
	}
	return &VarStmt{token, initilizer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(l.PRINT) != nil {
		return p.printStmt()
	}
	if p.match(l.LEFT_BRACE) != nil {
		return p.blockStmt()
	}
	if p.match(l.IF) != nil {
		return p.ifStmt()
	}
	if p.match(l.WHILE) != nil {
		return p.whileStmt()
	}
	if p.match(l.FOR) != nil {
		return p.forStmt()
	}
	if p.match(l.BREAK) != nil {
		return p.breakStmt()
	}
	if p.match(l.CONTINUE) != nil {
		return p.continueStmt()
	}
	if p.match(l.RETURN) != nil {
		return p.returnStmt()
	}
	return p.exprStmt()
}

func (p *Parser) returnStmt() (Stmt, error) {
	returnToken := p.previous()
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(l.SEMICOLON, "Expected ';' after return"); err != nil {
		return nil, err
	}
	return &ReturnStmt{returnToken, expr}, nil
}

func (p *Parser) continueStmt() (Stmt, error) {
	if err := p.consume(l.SEMICOLON, "Expected ';' after continue"); err != nil {
		return nil, err
	}
	return &ContinueStmt{p.previous()}, nil
}

func (p *Parser) breakStmt() (Stmt, error) {
	if err := p.consume(l.SEMICOLON, "Expected ';' after break"); err != nil {
		return nil, err
	}
	return &BreakStmt{p.previous()}, nil
}

func (p *Parser) forStmt() (Stmt, error) {
	if err := p.consume(l.LEFT_PAREN, "Expected '(' after for"); err != nil {
		return nil, err
	}
	var initialization Stmt = nil
	var err error = nil
	if p.match(l.VAR) != nil {
		initialization, err = p.varDecl()
	} else {
		initialization, err = p.exprStmt()
	}
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(l.SEMICOLON, "Expected ';' after condition"); err != nil {
		return nil, err
	}
	updation, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(l.RIGHT_PAREN, "Expected ')' after updation"); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &ForStmt{
		Initialization: initialization,
		Condition:      condition,
		Updation:       updation,
		Body:           body,
	}, nil
}

func (p *Parser) whileStmt() (Stmt, error) {
	if err := p.consume(l.LEFT_PAREN, "Expected '(' after while"); err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(l.RIGHT_PAREN, "Expected ')' after condition"); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &WhileStmt{
		Condition: expr,
		Body:      body,
	}, nil
}

func (p *Parser) ifStmt() (Stmt, error) {
	if err := p.consume(l.LEFT_PAREN, "Expected '(' after if"); err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(l.RIGHT_PAREN, "Expected ')' after condition"); err != nil {
		return nil, err
	}
	thenStmt, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseStmt Stmt = nil
	if p.match(l.ELSE) != nil {
		elseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &IfStmt{
		Condition: expr,
		ThenStmt:  thenStmt,
		ElseStmt:  elseStmt,
	}, nil
}

func (p *Parser) blockStmt() (Stmt, error) {
	declarations := make([]Stmt, 0)
	for !p.isAtEnd() && p.peek().Type != l.RIGHT_BRACE {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		declarations = append(declarations, declaration)
	}
	if err := p.consume(l.RIGHT_BRACE, "Expected '}' after block"); err != nil {
		return nil, err
	}
	return &BlockStmt{
		Declarations: declarations,
	}, nil
}

func (p *Parser) printStmt() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err = p.consume(l.SEMICOLON, "Expected ';' after value"); err != nil {
		return nil, err
	}
	return &PrintStmt{expr}, nil
}

func (p *Parser) exprStmt() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, ok := expr.(*FuncExpr); !ok {
		if err = p.consume(l.SEMICOLON, "Expected ';' after expression"); err != nil {
			return nil, err
		}
	}
	if expr == nil {
		return nil, nil
	}
	return &ExprStmt{expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.binaryPrec(binRules, LOGICAL_OR)
	if err != nil {
		return nil, err
	}

	if equal := p.match(l.EQUAL); equal != nil {
		expr, ok := expr.(*VariableExpr)
		if !ok {
			return nil, NewParserError(equal, "Invalid assignment target.")
		}
		assignment, err := p.assignment()
		if err != nil {
			return nil, err
		}
		return &AssignExpr{expr.Name, assignment}, nil
	}

	return expr, nil
}

// `rules` parameter is an array of BinaryRule that were defined
// with the priority go from highest to lowest accoding to its index.
// `binaryPrec` should be passed zero value for `ruleIndex` parameter
func (p *Parser) binaryPrec(rules []BinaryRule, ruleIndex int) (Expr, error) {
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

		expr = &BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    childPrec,
		}
	}
}

func (p *Parser) unary() (Expr, error) {
	operator := p.match(l.BANG, l.MINUS)
	if operator == nil {
		return p.call()
	}
	unaryExpr, err := p.unary()
	if err != nil {
		return nil, err
	}
	return &UnaryExpr{operator, unaryExpr}, nil
}
func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(l.LEFT_PAREN) != nil {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)
	if p.peek().Type != l.RIGHT_PAREN {
		for {
			argument, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argument)
			if p.match(l.COMMA) == nil {
				break
			}
		}
	}

	err := p.consume(l.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &CallExpr{
		Callee:    callee,
		Arguments: arguments,
	}, nil
}

func (p *Parser) primary() (Expr, error) {
	token := p.advance()
	switch token.Type {
	case l.NUMBER, l.STRING, l.TRUE, l.FALSE, l.NIL:
		return &PrimaryExpr{token}, nil
	case l.LEFT_PAREN:
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if err = p.consume(l.RIGHT_PAREN, "Expected ')' after expression"); err != nil {
			return nil, err
		}
		return expr, nil
	case l.IDENTIFIER:
		return &VariableExpr{token}, nil
	case l.FUN:
		return p.function()
	default:
		// Give back the token if don't match any precedence
		p.current--
		return nil, nil
	}
}

func (p *Parser) function() (Expr, error) {
	funcName := p.match(l.IDENTIFIER)

	err := p.consume(l.LEFT_PAREN, "Expect '(' after function name.")
	if err != nil {
		return nil, err
	}
	parameters := make([]*l.Token, 0)
	if p.peek().Type != l.RIGHT_PAREN {
		for {
			err := p.consume(l.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, p.previous())

			if p.match(l.COMMA) != nil {
				continue
			} else {
				break
			}
		}
	}
	err = p.consume(l.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}
	err = p.consume(l.LEFT_BRACE, "Expect '{' after function declaration.")
	if err != nil {
		return nil, err
	}
	body, err := p.blockStmt()
	if err != nil {
		return nil, err
	}
	return &FuncExpr{
		FuncStmt: &FuncStmt{
			Name:       funcName,
			Parameters: parameters,
			Body:       body.(*BlockStmt),
		},
	}, nil
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.current].Type == l.EOF
}

func (p *Parser) advance() *l.Token {
	if p.isAtEnd() {
		return nil
	}
	token := p.tokens[p.current]
	p.current++
	return &token
}

func (p *Parser) previous() *l.Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) peek() *l.Token {
	if p.isAtEnd() {
		return nil
	}
	return &p.tokens[p.current]
}

func (p *Parser) match(types ...l.TokenType) *l.Token {
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
		if previous.Type == l.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case l.CLASS, l.FUN, l.VAR, l.FOR, l.IF, l.WHILE, l.PRINT, l.RETURN:
			return
		}

		previous = p.advance()
	}
}

func (p *Parser) consume(tokenType l.TokenType, message string) error {
	if p.tokens[p.current].Type != tokenType {
		return NewParserError(&p.tokens[p.current], message)
	}
	p.current++
	return nil
}
