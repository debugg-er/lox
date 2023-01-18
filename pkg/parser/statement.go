package parser

import (
	"fmt"

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
	Execute()
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

func (t *PrintStmt) Execute() {
	value, _ := t.Expr.Evaluate()
	fmt.Print(value.Data)
}

func (t *ExprStmt) Execute() {
	t.Expr.Display(0)
}

func (t *VarStmt) Execute() {

}
