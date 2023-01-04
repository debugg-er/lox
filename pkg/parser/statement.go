package parser

import (
	"fmt"
)

type StmtType int

const (
	PRINT_STMM = iota
	EXPR_STMT
	IF_STMT
	FOR_STMT
	WHILE_STMT
	RETURN_STMT
)

type Stmt struct {
	Type StmtType
	Expr *Expr
}

func (t *Stmt) Execute() {
	switch t.Type {
	case EXPR_STMT:
		t.Expr.Display(0)
	case PRINT_STMM:
		value, _ := t.Expr.Evaluate()
		fmt.Print(value.Data)
	}
}
