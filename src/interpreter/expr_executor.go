package interpreter

import (
	l "github.com/debugg-er/lox/src/lexer"
	"github.com/debugg-er/lox/src/parser"
)

func (i *Interpreter) Evaluate(e parser.Expr) (*Value, error) {
	switch e := e.(type) {
	case *parser.PrimaryExpr:
		return i.evaluatePrimary(e)
	case *parser.UnaryExpr:
		return i.evaluateUnary(e)
	case *parser.BinaryExpr:
		return i.evaluateBinary(e)
	case *parser.VariableExpr:
		return i.evaluateVariable(e)
	case *parser.AssignExpr:
		return i.evaluateAssign(e)
	case *parser.FuncExpr:
		return i.evaluateFunc(e)
	case *parser.CallExpr:
		return i.evaluateCall(e)
	}

	return nil, nil
}

func (i *Interpreter) evaluatePrimary(e *parser.PrimaryExpr) (*Value, error) {
	return NewValue(e.Value.Value), nil
}

func (i *Interpreter) evaluateUnary(e *parser.UnaryExpr) (*Value, error) {
	preValue, err := i.Evaluate(e)
	if err != nil {
		return nil, err
	}
	switch e.Operator.Type {
	case l.MINUS:
		if isNumericOperand(*preValue) {
			return NewValue(-toNumber(*preValue)), nil
		} else {
			return nil, NewRuntimeError(e.Operator, "Bad datatype for unary operator")
		}
	case l.BANG:
		return NewValue(!isTruthy(*preValue)), nil
	default:
		panic("Language fatal: Undefined unary operator")
	}
}

func (i *Interpreter) evaluateBinary(e *parser.BinaryExpr) (*Value, error) {
	switch e.Operator.Type {
	case l.AND, l.OR:
		left, err := i.Evaluate(e.Left)
		if err != nil {
			return nil, err
		}
		if isTruthy(*left) && e.Operator.Type == l.OR {
			return NewValue(true), nil
		}
		right, err := i.Evaluate(e.Right)
		if err != nil {
			return nil, err
		}
		return NewValue(isTruthy(*right)), nil
	}

	left, err := i.Evaluate(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Type {
	case l.PLUS:
		if !isNumericOperand(*left, *right) {
			return NewValue(left.Stringify() + right.Stringify()), nil
		}
		return NewValue(toNumber(*left) + toNumber(*right)), nil
	case l.MINUS:
		if !isNumericOperand(*left, *right) {
			return nil, NewRuntimeError(e.Operator, "Operands must be a number")
		}
		return NewValue(toNumber(*left) - toNumber(*right)), nil
	case l.STAR:
		if !isNumericOperand(*left, *right) {
			return nil, NewRuntimeError(e.Operator, "Operands must be a number")
		}
		return NewValue(toNumber(*left) * toNumber(*right)), nil
	case l.SLASH:
		if !isNumericOperand(*left, *right) {
			return nil, NewRuntimeError(e.Operator, "Operands must be a number")
		}
		if toNumber(*right) == 0 {
			return nil, NewRuntimeError(e.Operator, "Division by zero")
		}
		return NewValue(toNumber(*left) / toNumber(*right)), nil
	case l.EQUAL_EQUAL:
		return NewValue(left.DataType == right.DataType && left.Data == right.Data), nil
	case l.BANG_EQUAL:
		return NewValue(!(left.DataType == right.DataType && left.Data == right.Data)), nil
	case l.GREATER_EQUAL:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) >= right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) >= toNumber(*right)), nil
		}
		return nil, NewRuntimeError(e.Operator, "Incompatible operands")
	case l.GREATER:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) > right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) > toNumber(*right)), nil
		}
		return nil, NewRuntimeError(e.Operator, "Incompatible operands")
	case l.LESS_EQUAL:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) <= right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) <= toNumber(*right)), nil
		}
		return nil, NewRuntimeError(e.Operator, "Incompatible operands")
	case l.LESS:
		if left.DataType == STRING_DT && right.DataType == STRING_DT {
			return NewValue(left.Data.(string) < right.Data.(string)), nil
		}
		if isNumericOperand(*left, *right) {
			return NewValue(toNumber(*left) < toNumber(*right)), nil
		}
		return nil, NewRuntimeError(e.Operator, "Incompatible operands")
	default:
		panic("Language fatal: Undefined binary operator")
	}

}

func (i *Interpreter) evaluateVariable(e *parser.VariableExpr) (*Value, error) {
	return i.env.get(e.Name)
}

func (i *Interpreter) evaluateAssign(e *parser.AssignExpr) (*Value, error) {
	value, err := i.Evaluate(e.Value)
	if err != nil {
		return nil, err
	}
	err = i.env.assign(e.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) evaluateFunc(e *parser.FuncExpr) (*Value, error) {
	value := &Value{
		DataType: FUNCTION_DT,
		Data:     e.FuncStmt,
	}
	if e.FuncStmt.Name != nil {
		i.env.define(e.FuncStmt.Name, value)
	}
	return value, nil
}

func (i *Interpreter) evaluateCall(e *parser.CallExpr) (*Value, error) {
	value, err := i.Evaluate(e.Callee)
	if err != nil {
		return nil, err
	}
	if value.DataType != FUNCTION_DT {
		// Must fix
		return nil, NewRuntimeError(&l.Token{Type: l.FUN, Line: 0}, "Expected function call.")
	}
	funcStmt := value.Data.(*parser.FuncStmt)
	if len(e.Arguments) < len(funcStmt.Parameters) {
		return nil, NewRuntimeError(funcStmt.Parameters[0], "Too few arguments.")
	}
	if len(e.Arguments) > len(funcStmt.Parameters) {
		return nil, NewRuntimeError(funcStmt.Parameters[0], "Too many arguments.")
	}

	oldEnv := i.env
	i.env = NewEnvironment(i.env)
	i.env.returableTarget = funcStmt
	defer func() {
		i.env = oldEnv
		funcStmt.SetIsReturned(false)
	}()

	for j, paramName := range funcStmt.Parameters {
		argumentVal, err := i.Evaluate(e.Arguments[j])
		if err != nil {
			return nil, err
		}
		i.env.define(paramName, argumentVal)
	}

	if err = i.Execute(funcStmt); err != nil {
		return nil, err
	}
	if i.env.returnValue == nil {
		return NewValue(nil), nil
	}
	return i.env.returnValue, nil
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
