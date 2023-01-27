package parser

import (
	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type Environment struct {
	store     map[string]*Value
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{make(map[string]*Value), enclosing}
}

func (e *Environment) define(variable *Token, value *Value) {
	e.store[variable.Value.(string)] = value
}

func (e *Environment) get(variable *Token) (*Value, *Error) {
	varName := variable.Value.(string)
	value := e.store[varName]
	if value != nil {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.get(variable)
	}
	return nil, NewError(variable, "Undefined variable '"+varName+"'.")
}

func (e *Environment) assign(variable *Token, value *Value) *Error {
	varName := variable.Value.(string)
	if e.store[varName] != nil {
		e.store[variable.Value.(string)] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.assign(variable, value)
	}
	return NewError(variable, "Undefined variable '"+varName+"'.")
}
