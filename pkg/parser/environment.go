package parser

import (
	. "github.com/debugg-er/lox/pkg/common"
)

type Environment struct {
	store map[string]*Value
}

func (e *Environment) define(name *Token, value *Value) {
	e.store[name.Value.(string)] = value
}

func (e *Environment) get(name *Token) (*Value, *Error) {
	value := e.store[name.Value.(string)]

	if value == nil {
		return nil, NewError(name, "Undefined variable '"+name.Value.(string)+"'.")
	}
	return value, nil
}

func (e *Environment) assign(name *Token, value *Value) *Error {
	if e.store[name.Value.(string)] == nil {
		return NewError(name, "Undefined variable '"+name.Value.(string)+"'.")
	}
	e.store[name.Value.(string)] = value
	return nil
}
