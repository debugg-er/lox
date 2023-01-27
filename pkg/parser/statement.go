package parser

import (
	"fmt"
	"strings"

	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type StmtType int

const (
	PRINT_STMM = iota
	EXPR_STMT
	IF_STMT
	FOR_STMT
	WHILE_STMT
	RETURN_STMT
	VAR_STMT
)

type Stmt interface {
	Execute(e *Environment) *Error
}

type PrintStmt struct {
	Expr *Expr
}

type ExprStmt struct {
	Expr *Expr
}

type VarStmt struct {
	name       *Token
	initilizer *Expr
}

type BlockStmt struct {
	declarations []Stmt
}

type IfStmt struct {
	condition *Expr
	trueStmt  Stmt
	falseStmt Stmt
}

func (t *PrintStmt) Execute(e *Environment) *Error {
	value, err := t.Expr.Evaluate(e)
	if err != nil {
		return err
	}
	fmt.Print(strings.ReplaceAll(value.Stringify(), `\n`, "\n"))
	return nil
}

func (t *ExprStmt) Execute(e *Environment) *Error {
	if t.Expr.Type == ASSIGN {
		t.Expr.Evaluate(e)
	} else {
		t.Expr.Display(0)
	}
	return nil
}

func (t *VarStmt) Execute(e *Environment) *Error {
	value, err := t.initilizer.Evaluate(e)
	if err != nil {
		return err
	}
	e.define(t.name, value)
	return nil
}

func (t *BlockStmt) Execute(e *Environment) *Error {
	blockEnv := NewEnvironment(e)
	for _, stmt := range t.declarations {
		err := stmt.Execute(blockEnv)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *IfStmt) Execute(e *Environment) *Error {
	conditionValue, err := t.condition.Evaluate(e)
	if err != nil {
		return err
	}
	if isTruthy(*conditionValue) {
		t.trueStmt.Execute(e)
	} else if t.falseStmt != nil {
		t.falseStmt.Execute(e)
	}
	return nil
}
