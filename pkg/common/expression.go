package common

import (
	"fmt"
	"strings"
)

type ExprType int

const (
	LITERAL = iota
	GROUPING
	UNARY
	BINARY
)

type Expr struct {
	Type     ExprType
	Operator TokenType
	Left     *Expr
	Right    *Expr
	Literal  interface{}
}

func DisplayExpr(expr *Expr, tab int) {
	if expr == nil {
		fmt.Println("nil")
		return
	}

	fmt.Println(strings.Repeat(" ", tab), "{")
	fmt.Println(strings.Repeat(" ", tab), "operator: ", expr.Operator)
	fmt.Print(strings.Repeat(" ", tab), "left: ")
	DisplayExpr(expr.Left, tab+4)
	fmt.Print(strings.Repeat(" ", tab), "right: ")
	DisplayExpr(expr.Right, tab+4)
	fmt.Println(strings.Repeat(" ", tab), "}")
}
