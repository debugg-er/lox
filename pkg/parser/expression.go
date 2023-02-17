package parser

import (
	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type (
	Expr interface {
		Evaluate(env *Environment) (*Value, error)
	}

	Callable interface {
		Call()
	}
)

type (
	PrimaryExpr struct {
		value *Token
	}

	UnaryExpr struct {
		operator *Token
		operand  Expr
	}

	BinaryExpr struct {
		operator *Token
		left     Expr
		right    Expr
	}

	VariableExpr struct {
		name *Token
	}

	AssignExpr struct {
		name  *Token
		value Expr
	}

	FuncExpr struct {
		funcStmt *FuncStmt
	}

	CallExpr struct {
		callee    Expr
		arguments []Expr
	}
)

func (e *PrimaryExpr) Evaluate(env *Environment) (*Value, error) {
	return NewValue(e.value.Value), nil
}

func (e *UnaryExpr) Evaluate(env *Environment) (*Value, error) {
	preValue, err := e.operand.Evaluate(env)
	if err != nil {
		return nil, err
	}
	switch e.operator.Type {
	case MINUS:
		if isNumericOperand(*preValue) {
			return NewValue(-toNumber(*preValue)), nil
		} else {
			return nil, NewError(e.operator, "Bad datatype for unary operator")
		}
	case BANG:
		return NewValue(!isTruthy(*preValue)), nil
	default:
		panic("Language fatal: Undefined unary operator")
	}
}

func (e *BinaryExpr) Evaluate(env *Environment) (*Value, error) {
	switch e.operator.Type {
	case AND, OR:
		left, err := e.left.Evaluate(env)
		if err != nil {
			return nil, err
		}
		if isTruthy(*left) && e.operator.Type == OR {
			return NewValue(true), nil
		}
		right, err := e.right.Evaluate(env)
		if err != nil {
			return nil, err
		}
		return NewValue(isTruthy(*right)), nil
	}

	left, err := e.left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := e.right.Evaluate(env)
	if err != nil {
		return nil, err
	}

	switch e.operator.Type {
	case PLUS:
		if !isNumericOperand(*left, *right) {
			return NewValue(left.Stringify() + right.Stringify()), nil
		}
		return NewValue(toNumber(*left) + toNumber(*right)), nil
	case MINUS:
		if !isNumericOperand(*left, *right) {
			return nil, NewError(e.operator, "Operands must be a number")
		}
		return NewValue(toNumber(*left) - toNumber(*right)), nil
	case STAR:
		if !isNumericOperand(*left, *right) {
			return nil, NewError(e.operator, "Operands must be a number")
		}
		return NewValue(toNumber(*left) * toNumber(*right)), nil
	case SLASH:
		if !isNumericOperand(*left, *right) {
			return nil, NewError(e.operator, "Operands must be a number")
		}
		if toNumber(*right) == 0 {
			return nil, NewError(e.operator, "Division by zero")
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
		return nil, NewError(e.operator, "Incompatible operands")
	case GREATER:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) > right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) > toNumber(*right)), nil
		}
		return nil, NewError(e.operator, "Incompatible operands")
	case LESS_EQUAL:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) <= right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) <= toNumber(*right)), nil
		}
		return nil, NewError(e.operator, "Incompatible operands")
	case LESS:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) < right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) < toNumber(*right)), nil
		}
		return nil, NewError(e.operator, "Incompatible operands")
	default:
		panic("Language fatal: Undefined binary operator")
	}

}

func (e *VariableExpr) Evaluate(env *Environment) (*Value, error) {
	return env.get(e.name)
}

func (e *AssignExpr) Evaluate(env *Environment) (*Value, error) {
	value, err := e.value.Evaluate(env)
	if err != nil {
		return nil, err
	}
	err = env.assign(e.name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (e *FuncExpr) Evaluate(env *Environment) (*Value, error) {
	value := &Value{
		DataType: FUNCTION_DT,
		Data:     e.funcStmt,
	}
	if e.funcStmt.name != nil {
		env.define(e.funcStmt.name, value)
	}
	return value, nil
}

func (e *CallExpr) Evaluate(env *Environment) (*Value, error) {
	value, err := e.callee.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if value.DataType != FUNCTION_DT {
		// Must fix
		return nil, NewError(&Token{Type: FUN, Line: 0}, "Expected function call.")
	}
	funcStmt := value.Data.(*FuncStmt)
	if len(e.arguments) < len(funcStmt.parameters) {
		return nil, NewError(funcStmt.parameters[0], "Too few arguments.")
	}
	if len(e.arguments) > len(funcStmt.parameters) {
		return nil, NewError(funcStmt.parameters[0], "Too many arguments.")
	}
	funcEnv := NewEnvironment(env)
	for i, token := range funcStmt.parameters {
		argumentVal, err := e.arguments[i].Evaluate(env)
		if err != nil {
			return nil, err
		}
		funcEnv.define(token, argumentVal)
	}

	err = funcStmt.Execute(funcEnv)
	if err != nil {
		return nil, err
	}
	return funcStmt.getReturnValue(), nil
}

// func (e *CallExpr) Call(env *Environment) (*Value, error) {

// }

func isTruthy(value Value) bool {
	switch value.DataType {
	case NUMBER_DT:
		return value.Data != 0
	case STRING_DT:
		return value.Data != ""
	case BOOLEAN_DT:
		return value.Data.(bool)
	case NULL_DT:
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
