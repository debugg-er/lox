package parser

import (
	"fmt"
	"strings"

	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type ExprType int

const (
	PRIMARY = iota
	// GROUPING
	UNARY
	BINARY
	VARIABLE
	ASSIGN
)

type Expr struct {
	Type        ExprType
	Operator    *Token
	Left        *Expr
	Right       *Expr
	Primary     *Token
	Var         *Token
	AssignValue *Expr
}

func (e *Expr) Evaluate(env *Environment) (*Value, *Error) {
	switch e.Type {
	case VARIABLE:
		return e.evaluateVariable(env)
	case ASSIGN:
		return e.evaluateAssignment(env)
	case PRIMARY:
		return e.evaluatePrimary(env)
	case UNARY:
		return e.evaluateUnary(env)
	case BINARY:
		return e.evaluateBinary(env)
	default:
		return nil, nil
	}
}

func NewPrimaryExpr(value *Token) *Expr {
	return &Expr{
		Type:    PRIMARY,
		Primary: value,
	}
}

func (e *Expr) evaluateVariable(env *Environment) (*Value, *Error) {
	return env.get(e.Var)
}

func (e *Expr) evaluateAssignment(env *Environment) (*Value, *Error) {
	value, err := e.AssignValue.Evaluate(env)
	if err != nil {
		return nil, err
	}
	err = env.assign(e.Var, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (e *Expr) evaluateBinary(env *Environment) (*Value, *Error) {
	switch e.Operator.Type {
	case AND, OR:
		return e.evaluateLogicalBinary(env)
	}
	left, err := e.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := e.Right.Evaluate(env)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Type {
	case PLUS:
		if !isNumericOperand(*left, *right) {
			return NewValue(left.Stringify() + right.Stringify()), nil
		}
		return NewValue(toNumber(*left) + toNumber(*right)), nil
	case MINUS:
		if !isNumericOperand(*left, *right) {
			return nil, NewError(e.Operator, "Operands must be a number")
		}
		return NewValue(toNumber(*left) - toNumber(*right)), nil
	case STAR:
		if !isNumericOperand(*left, *right) {
			return nil, NewError(e.Operator, "Operands must be a number")
		}
		return NewValue(toNumber(*left) * toNumber(*right)), nil
	case SLASH:
		if !isNumericOperand(*left, *right) {
			return nil, NewError(e.Operator, "Operands must be a number")
		}
		if toNumber(*right) == 0 {
			return nil, NewError(e.Operator, "Division by zero")
		}
		return NewValue(toNumber(*left) / toNumber(*right)), nil
	case EQUAL_EQUAL:
		return NewValue(left.DataType == right.DataType && left.Data == right.Data), nil
	case BANG_EQUAL:
		return NewValue(!(left.DataType == right.DataType && left.Data == right.Data)), nil
	case GREATER_EQUAL:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) >= right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) >= toNumber(*right)), nil
		}
		return nil, NewError(e.Operator, "Incompatible operands")
	case GREATER:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) > right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) > toNumber(*right)), nil
		}
		return nil, NewError(e.Operator, "Incompatible operands")
	case LESS_EQUAL:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) <= right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) <= toNumber(*right)), nil
		}
		return nil, NewError(e.Operator, "Incompatible operands")
	case LESS:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) < right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) < toNumber(*right)), nil
		}
		return nil, NewError(e.Operator, "Incompatible operands")
	default:
		panic("Language fatal: Undefined binary operator")
	}
}

func (e *Expr) evaluateLogicalBinary(env *Environment) (*Value, *Error) {
	switch e.Operator.Type {
	case AND:
		left, err := e.Left.Evaluate(env)
		if err != nil {
			return nil, err
		}
		if !isTruthy(*left) {
			return NewValue(false), nil
		}
		right, err := e.Right.Evaluate(env)
		if err != nil {
			return nil, err
		}
		return NewValue(isTruthy(*right)), nil
	case OR:
		left, err := e.Left.Evaluate(env)
		if err != nil {
			return nil, err
		}
		if isTruthy(*left) {
			return NewValue(true), nil
		}
		right, err := e.Right.Evaluate(env)
		if err != nil {
			return nil, err
		}
		return NewValue(isTruthy(*right)), nil
	default:
		return nil, nil
	}
}

func (e *Expr) evaluateUnary(env *Environment) (*Value, *Error) {
	preValue, err := e.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	switch e.Operator.Type {
	case MINUS:
		if isNumericOperand(*preValue) {
			return NewValue(-toNumber(*preValue)), nil
		} else {
			return nil, NewError(e.Operator, "Bad datatype for unary operator")
		}
	case BANG:
		return NewValue(!isTruthy(*preValue)), nil
	default:
		panic("Language fatal: Undefined unary operator")
	}
}

func (e *Expr) evaluatePrimary(env *Environment) (*Value, *Error) {
	return NewValue(e.Primary.Value), nil
}

func (e *Expr) Display(tab int) {
	if e == nil {
		fmt.Println("nil")
		return
	}

	spaces := strings.Repeat(" ", tab)
	spaces4 := spaces + "    "

	fmt.Println("{")
	if e.Type == VARIABLE {
		fmt.Printf("%s%s%v\n", spaces4, "variable: ", e.Var.Value)
	} else if e.Type == ASSIGN {
		fmt.Printf("%s%s%v", spaces4, "variable: ", e.Var.Value)
		e.Left.Display(tab + 4)
	} else if e.Operator != nil {
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
		panic("Language fatal: Undefined datatype")
	}
}

func isNumericOperand(values ...Value) bool {
	for _, value := range values {
		if value.DataType != BOOLEAN_DT && value.DataType != NUMBER_DT {
			return false
		}
	}
	return true
}

func toNumber(value Value) float64 {
	switch value.DataType {
	case BOOLEAN_DT:
		if value.Data.(bool) {
			return 1
		} else {
			return 0
		}
	case NUMBER_DT:
		return value.Data.(float64)
	default:
		panic("Language Fatal: Can't parse value to number")
	}
}
