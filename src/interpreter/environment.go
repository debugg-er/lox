package interpreter

import (
	l "github.com/debugg-er/lox/src/lexer"
	"github.com/debugg-er/lox/src/parser"
)

type Environment struct {
	store           map[string]*Value
	enclosing       *Environment
	loopableTarget  parser.Loopable   // Executor Target (ForStmt, WhileStmt)
	returableTarget parser.Returnable // Executor Target (FuncStmt)
	returnValue     *Value
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		store:          make(map[string]*Value),
		enclosing:      enclosing,
		loopableTarget: nil,
	}
}

func (e *Environment) define(variable *l.Token, value *Value) {
	e.store[variable.Value.(string)] = value
}

func (e *Environment) get(variable *l.Token) (*Value, error) {
	varName := variable.Value.(string)
	value := e.store[varName]
	if value != nil {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.get(variable)
	}
	return nil, NewRuntimeError(variable, "Undefined variable '"+varName+"'.")
}

func (e *Environment) assign(variable *l.Token, value *Value) error {
	varName := variable.Value.(string)
	if e.store[varName] != nil {
		e.store[variable.Value.(string)] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.assign(variable, value)
	}
	return NewRuntimeError(variable, "Undefined variable '"+varName+"'.")
}

func (e *Environment) getLoopableTarget() parser.Loopable {
	if e.loopableTarget != nil {
		return e.loopableTarget
	}
	if e.enclosing != nil {
		return e.enclosing.getLoopableTarget()
	}
	return nil
}

func (e *Environment) getReturnableTarget() parser.Returnable {
	if e.returableTarget != nil {
		return e.returableTarget
	}
	if e.enclosing != nil {
		return e.enclosing.getReturnableTarget()
	}
	return nil
}

func (e *Environment) getReturnableTargetEnv() *Environment {
	if e.returableTarget != nil {
		return e
	}
	if e.enclosing != nil {
		return e.enclosing.getReturnableTargetEnv()
	}
	return nil
}
