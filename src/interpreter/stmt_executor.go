package interpreter

import (
	"fmt"

	"github.com/debugg-er/lox/src/parser"
)

func (i *Interpreter) Execute(t parser.Stmt) error {
	switch t := t.(type) {
	case *parser.PrintStmt:
		return i.executePrintStmt(t)
	case *parser.ExprStmt:
		return i.executeExprStmt(t)
	case *parser.VarStmt:
		return i.executeVarStmt(t)
	case *parser.BlockStmt:
		return i.executeBlockStmt(t)
	case *parser.IfStmt:
		return i.executeIfStmt(t)
	case *parser.WhileStmt:
		return i.executeWhileStmt(t)
	case *parser.ForStmt:
		return i.executeForStmt(t)
	case *parser.BreakStmt:
		return i.executeBreakStmt(t)
	case *parser.ContinueStmt:
		return i.executeContinueStmt(t)
	case *parser.FuncStmt:
		return i.executeFuncStmt(t)
	case *parser.ReturnStmt:
		return i.executeReturnStmt(t)
	}
	return nil
}

// ---------------- Print Statement ----------------
func (i *Interpreter) executePrintStmt(t *parser.PrintStmt) error {
	value, err := i.Evaluate(t.Expr)
	if err != nil {
		return err
	}
	fmt.Println(value.Stringify())
	return nil
}

// ---------------- Expression Statement ----------------
func (i *Interpreter) executeExprStmt(t *parser.ExprStmt) error {
	_, err := i.Evaluate(t.Expr)
	return err
}

// ---------------- Variable Declaration Statement ----------------
func (i *Interpreter) executeVarStmt(t *parser.VarStmt) error {
	value, err := i.Evaluate(t.Initilizer)
	if err != nil {
		return err
	}
	i.env.define(t.Name, value)
	return nil
}

// ---------------- Block Statement ----------------
func (i *Interpreter) executeBlockStmt(t *parser.BlockStmt) error {
	oldEnv := i.env
	i.env = NewEnvironment(i.env)
	defer func(env *Environment) {
		i.env = env
	}(oldEnv)

	for _, stmt := range t.Declarations {
		err := i.Execute(stmt)
		if err != nil {
			return err
		}
		// Stop execute when break or continue is met on child statements
		if executor := i.env.getLoopableTarget(); executor != nil {
			if executor.IsBreaked() || executor.IsContinued() {
				return nil
			}
		}
		// Stop execute when return is met on child statements
		if executor := i.env.getReturnableTarget(); executor != nil {
			if executor.IsReturned() {
				return nil
			}
		}
	}
	return nil
}

// ---------------- If Statement ----------------
func (i *Interpreter) executeIfStmt(t *parser.IfStmt) error {
	conditionValue, err := i.Evaluate(t.Condition)
	if err != nil {
		return err
	}
	if isTruthy(*conditionValue) {
		i.Execute(t.ThenStmt)
	} else if t.ElseStmt != nil {
		i.Execute(t.ElseStmt)
	}
	return nil
}

// ---------------- While Statement ----------------
func (i *Interpreter) executeWhileStmt(t *parser.WhileStmt) error {
	i.env.loopableTarget = t
	for {
		conditionValue, err := i.Evaluate(t.Condition)
		if err != nil {
			return err
		}
		if !isTruthy(*conditionValue) {
			return nil
		}
		i.Execute(t.Body)
		if t.IsBreaked() {
			return nil
		}
		if t.IsContinued() {
			t.SetContinued(false)
		}
	}
}

// ---------------- For Statement ----------------
func (i *Interpreter) executeForStmt(t *parser.ForStmt) error {
	i.env.loopableTarget = t
	if t.Initialization != nil {
		i.Execute(t.Initialization)
	}
	for {
		// Condition checking
		if t.Condition != nil {
			conditionValue, err := i.Evaluate(t.Condition)
			if err != nil {
				return err
			}
			if !isTruthy(*conditionValue) {
				return nil
			}
		}
		// Body execution
		i.Execute(t.Body)
		if t.IsBreaked() {
			return nil
		}
		if t.IsContinued() {
			t.SetContinued(false)
		}
		if t.Updation != nil {
			i.Evaluate(t.Updation)
		}
	}
}

// ---------------- Break Statement ----------------
func (i *Interpreter) executeBreakStmt(t *parser.BreakStmt) error {
	executor := i.env.getLoopableTarget()
	if executor != nil {
		return NewRuntimeError(t.Token, "RuntimeError: 'break' statement can only be used within an enclosing iteration")
	}
	executor.SetBreaked(true)
	return nil
}

// ---------------- Continue Statement ----------------
func (i *Interpreter) executeContinueStmt(t *parser.ContinueStmt) error {
	executor := i.env.getLoopableTarget()
	if executor == nil {
		return NewRuntimeError(t.Token, "RuntimeError: 'continue' statement can only be used within an enclosing iteration")
	}
	executor.SetContinued(true)
	return nil
}

// ---------------- Function Statement ----------------
func (i *Interpreter) executeFuncStmt(t *parser.FuncStmt) error {
	if err := i.Execute(t.Body); err != nil {
		return err
	}
	return nil
}

// ---------------- Return Statement ----------------
func (i *Interpreter) executeReturnStmt(t *parser.ReturnStmt) error {
	executor := i.env.getReturnableTarget()
	if executor == nil {
		return NewRuntimeError(t.Token, "RuntimeError: 'return' statement can only be used within an function")
	}
	value, err := i.Evaluate(t.Expr)
	if err != nil {
		return err
	}
	executor.SetIsReturned(true)
	returnableTargetEnv := i.env.getReturnableTargetEnv()
	returnableTargetEnv.returnValue = value
	return nil
}
