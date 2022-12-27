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

func NewLiteralExpr(value interface{}) *Expr {
	return &Expr{
		Type:    LITERAL,
		Literal: value,
	}
}

func (e *Expr) Display(tab int) {
	if e == nil {
		fmt.Println("nil")
		return
	}

	spaces := strings.Repeat(" ", tab)
	spaces4 := spaces + "    "

	fmt.Println("{")
	if e.Operator != Undefined {
		fmt.Printf("%s%s%v\n", spaces4, "operator: ", e.Operator)
	} else {
		fmt.Printf("%s%s%v\n", spaces4, "value: ", e.Literal)
	}
	if e.Left != nil {
		fmt.Print(spaces4 + "left: ")
		e.Left.Display(tab + 4)
	}
	if e.Right != nil {
		fmt.Print(spaces4 + "right: ")
		e.Right.Display(tab + 4)
	}
	fmt.Println(spaces + "}")
}
