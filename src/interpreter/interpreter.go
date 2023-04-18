package interpreter

import "github.com/debugg-er/lox/src/parser"

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(nil),
	}
}

func (i *Interpreter) Run(statements []parser.Stmt) error {
	for _, stmt := range statements {
		if err := i.Execute(stmt); err != nil {
			return err
		}
	}
	return nil
}
