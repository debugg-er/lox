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

func (e *Expr) Evaluate() (*Value, *Error) {
	switch e.Type {
	case PRIMARY:
		return evaluatePrimary(e)
	case UNARY:
		return evaluateUnary(e)
	case BINARY:
		return evaluateBinary(e)
	default:
		return nil, nil
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

func evaluateBinary(e *Expr) (*Value, *Error) {
	left, err := e.Left.Evaluate()
	if err != nil {
		return nil, err
	}
	right, err := e.Right.Evaluate()
	if err != nil {
		return nil, err
	}
	isBothNum := left.DataType == NUMBER_DT && right.DataType == NUMBER_DT
	l, _ := left.Data.(float64)
	r, _ := right.Data.(float64)

	switch e.Operator.Type {
	case PLUS:
		if !isBothNum {
			return NewValue(STRING_DT, left.Stringify()+right.Stringify()), nil
		}
		return NewValue(NUMBER_DT, l+r), nil
	case MINUS:
		if !isBothNum {
			return nil, NewError(e.Operator, "Operands must be a number")
		}
		return NewValue(NUMBER_DT, l-r), nil
	case STAR:
		if !isBothNum {
			return nil, NewError(e.Operator, "Operands must be a number")
		}
		return NewValue(NUMBER_DT, l*r), nil
	case SLASH:
		if !isBothNum {
			return nil, NewError(e.Operator, "Operands must be a number")
		}
		if r == 0 {
			return nil, NewError(e.Operator, "Division by zero")
		}
		return NewValue(NUMBER_DT, l/r), nil
	default:
		panic("Language fatal")
	}
}

func evaluateUnary(e *Expr) (*Value, *Error) {
	preValue, err := e.Left.Evaluate()
	if err != nil {
		return nil, err
	}
	switch e.Operator.Type {
	case MINUS:
		if preValue.DataType == NUMBER_DT {
			return NewValue(NUMBER_DT, -preValue.DataType), nil
		} else {
			return nil, NewError(e.Left.Primary, "Bad datatype for unary operator")
		}
	case BANG:
		return NewValue(BOOLEAN_DT, !isTruthy(*preValue)), nil
	default:
		panic("Language fatal")
	}
}

func evaluatePrimary(e *Expr) (*Value, *Error) {
	if e.Primary.Type == TRUE {
		return NewValue(NIL_DT, TRUE), nil
	}
	if e.Primary.Type == FALSE {
		return NewValue(NIL_DT, FALSE), nil
	}

	switch value := e.Primary.Value.(type) {
	case string:
		return NewValue(STRING_DT, value), nil
	case float64:
		return NewValue(NUMBER_DT, value), nil
	case nil:
		return NewValue(NIL_DT, value), nil
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
