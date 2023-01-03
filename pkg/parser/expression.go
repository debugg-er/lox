package parser

import (
	"fmt"
	"strings"

	. "github.com/debugg-er/lox/pkg/common"
)

type ExprType int

const (
	PRIMARY = iota
	// GROUPING
	UNARY
	BINARY
)

type Expr struct {
	Type     ExprType
	Operator *Token
	Left     *Expr
	Right    *Expr
	Primary  *Token
}

func (e *Expr) Evaluate() *Value {
	switch e.Type {
	case PRIMARY:
		return evaluatePrimary(e)
	case UNARY:
		return evaluateUnary(e)
	case BINARY:
		return evaluateBinary(e)
	default:
		return nil
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
	if e.Operator != nil {
		fmt.Printf("%s%s%v\n", spaces4, "operator: ", e.Operator.Type)
	} else {
		fmt.Printf("%s%s%v\n", spaces4, "value: ", e.Primary.Value)
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

func NewLiteralExpr(value *Token) *Expr {
	return &Expr{
		Type:    PRIMARY,
		Primary: value,
	}
}

func evaluateBinary(e *Expr) *Value {
	left := e.Left.Evaluate()
	right := e.Right.Evaluate()
	isBothNum := left.DataType == NUMBER_DT && right.DataType == NUMBER_DT

	switch e.Operator.Type {
	case PLUS:
		if !isBothNum {
			return NewValue(STRING_DT, left.Stringify()+right.Stringify())
		}
		return NewValue(NUMBER_DT, left.Data.(float64)+right.Data.(float64))
	case MINUS:
		if !isBothNum {
			// handle error
			fmt.Println("Can't substract non-number expression")
		}
		return NewValue(NUMBER_DT, left.Data.(float64)-right.Data.(float64))
	case STAR:
		if !isBothNum {
			// handle error
			fmt.Println("Can't multiply non-number expression")
		}
		return NewValue(NUMBER_DT, left.Data.(float64)*right.Data.(float64))
	case SLASH:
		if !isBothNum {
			// handle error
			fmt.Println("Can't divide non-number expression")
		}
		if right.Data.(float64) == 0 {
			// handle error
			fmt.Println("Can't divide with 0")
		}
		return NewValue(NUMBER_DT, left.Data.(float64)/right.Data.(float64))
	default:
		panic("Language fatal")
	}
}

func evaluateUnary(e *Expr) *Value {
	switch e.Operator.Type {
	case MINUS:
		preValue := e.Left.Evaluate()
		if preValue.DataType == NUMBER_DT {
			return NewValue(NUMBER_DT, -preValue.DataType)
		} else {
			// handle error
			fmt.Println("Expected number")
			return nil
		}
	case BANG:
		return NewValue(BOOLEAN_DT, !isTruthy(*e.Left.Evaluate()))
	default:
		panic("Language fatal")
	}
}

func evaluatePrimary(e *Expr) *Value {
	if e.Primary.Type == TRUE {
		return NewValue(NIL_DT, TRUE)
	}
	if e.Primary.Type == FALSE {
		return NewValue(NIL_DT, FALSE)
	}

	switch value := e.Primary.Value.(type) {
	case string:
		return NewValue(STRING_DT, value)
	case float64:
		return NewValue(NUMBER_DT, value)
	case nil:
		return NewValue(NIL_DT, value)
	default:
		panic("Language fatal: Undefined datatype")
	}
}

func isTruthy(value Value) bool {
	switch value.DataType {
	case NUMBER_DT:
		return value.Data != 0
	case STRING_DT:
		return value.Data != ""
	case BOOLEAN_DT:
		return value.Data.(bool)
	case NIL_DT:
		return false
	default:
		panic("Language fatal")
	}
}
