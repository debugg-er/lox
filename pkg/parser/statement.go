package parser

import (
	"fmt"

	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type (
	Stmt interface {
		Execute(e *Environment) *Error
	}

	Loopable interface {
		setBreaked(isBreaked bool)
		isBreaked() bool
		setContinued(isContinued bool)
		isContinued() bool
	}

	Returnable interface {
		setReturnValue(value *Value)
		getReturnValue() *Value
	}
)

type (
	PrintStmt struct {
		Expr *Expr
	}

	ExprStmt struct {
		Expr *Expr
	}

	VarStmt struct {
		name       *Token
		initilizer *Expr
	}

	BlockStmt struct {
		declarations []Stmt
	}

	IfStmt struct {
		condition *Expr
		thenStmt  Stmt
		elseStmt  Stmt
	}

	WhileStmt struct {
		condition    *Expr
		body         Stmt
		_isBreaked   bool
		_isContinued bool
	}

	ForStmt struct {
		initialization Stmt
		condition      *Expr
		updation       *Expr
		body           Stmt
		_isBreaked     bool
		_isContinued   bool
	}

	BreakStmt struct {
		token *Token
	}

	ContinueStmt struct {
		token *Token
	}

	FuncStmt struct {
		name        *Token
		arguments   any // placeholder
		body        *BlockStmt
		returnValue *Value
	}

	ReturnStmt struct {
		token *Token
		expr  *Expr
	}
)

// ---------------- Print Statement ----------------
func (t *PrintStmt) Execute(e *Environment) *Error {
	value, err := t.Expr.Evaluate(e)
	if err != nil {
		return err
	}
	fmt.Print(value.Stringify())
	return nil
}

// ---------------- Expression Statement ----------------
func (t *ExprStmt) Execute(e *Environment) *Error {
	if t.Expr.Type == ASSIGN {
		t.Expr.Evaluate(e)
	} else {
		t.Expr.Display(0)
	}
	return nil
}

// ---------------- Variable Declaration Statement ----------------
func (t *VarStmt) Execute(e *Environment) *Error {
	value, err := t.initilizer.Evaluate(e)
	if err != nil {
		return err
	}
	e.define(t.name, value)
	return nil
}

// ---------------- Block Statement ----------------
func (t *BlockStmt) Execute(e *Environment) *Error {
	blockEnv := NewEnvironment(e)
	for _, stmt := range t.declarations {
		err := stmt.Execute(blockEnv)
		if err != nil {
			return err
		}
		// Stop execute when break or continue is met on child statements
		if executor := e.getLoopableTarget(); executor != nil {
			if executor.isBreaked() || executor.isContinued() {
				return nil
			}
		}
		// Stop execute when break or continue is met
		switch stmt.(type) {
		case *BreakStmt, *ContinueStmt:
			return nil
		}
	}
	return nil
}

// ---------------- If Statement ----------------
func (t *IfStmt) Execute(e *Environment) *Error {
	conditionValue, err := t.condition.Evaluate(e)
	if err != nil {
		return err
	}
	if isTruthy(*conditionValue) {
		t.thenStmt.Execute(e)
	} else if t.elseStmt != nil {
		t.elseStmt.Execute(e)
	}
	return nil
}

// ---------------- While Statement ----------------
func (t *WhileStmt) Execute(e *Environment) *Error {
	e.loopableTarget = t
	for {
		conditionValue, err := t.condition.Evaluate(e)
		if err != nil {
			return err
		}
		if !isTruthy(*conditionValue) {
			return nil
		}
		t.body.Execute(e)
		if t.isContinued() {
			t.setContinued(false)
			continue
		}
		if t.isBreaked() {
			return nil
		}
	}
}

func (t *WhileStmt) isBreaked() bool {
	return t._isBreaked
}

func (t *WhileStmt) setBreaked(isBreaked bool) {
	t._isBreaked = isBreaked
}

func (t *WhileStmt) isContinued() bool {
	return t._isContinued
}

func (t *WhileStmt) setContinued(isContinued bool) {
	t._isContinued = isContinued
}

// ---------------- For Statement ----------------
func (t *ForStmt) Execute(e *Environment) *Error {
	e.loopableTarget = t
	if t.initialization != nil {
		t.initialization.Execute(e)
	}
	for {
		// Condition checking
		if t.condition != nil {
			conditionValue, err := t.condition.Evaluate(e)
			if err != nil {
				return err
			}
			if !isTruthy(*conditionValue) {
				return nil
			}
		}
		// Body execution
		t.body.Execute(e)
		if t.isContinued() {
			t.setContinued(false)
			continue
		}
		if t.isBreaked() {
			return nil
		}
		if t.updation != nil {
			t.updation.Evaluate(e)
		}
	}
}

func (t *ForStmt) isBreaked() bool {
	return t._isBreaked
}

func (t *ForStmt) setBreaked(isBreaked bool) {
	t._isBreaked = isBreaked
}

func (t *ForStmt) isContinued() bool {
	return t._isContinued
}

func (t *ForStmt) setContinued(isContinued bool) {
	t._isContinued = isContinued
}

// ---------------- Break Statement ----------------
func (t *BreakStmt) Execute(e *Environment) *Error {
	if executor := e.getLoopableTarget(); executor != nil {
		executor.setBreaked(true)
		return nil
	}
	return NewError(t.token, "RuntimeError: 'break' statement can only be used within an enclosing iteration")
}

// ---------------- Continue Statement ----------------
func (t *ContinueStmt) Execute(e *Environment) *Error {
	if executor := e.getLoopableTarget(); executor != nil {
		executor.setContinued(true)
		return nil
	}
	return NewError(t.token, "RuntimeError: 'continue' statement can only be used within an enclosing iteration")
}

// ---------------- Function Statement ----------------
func (t *FuncStmt) Execute(e *Environment) *Error {
	e.returableTarget = t
	// Todo: Implement function execution
	return nil
}

func (t *FuncStmt) setReturnValue(value *Value) {
	t.returnValue = value
}

func (t *FuncStmt) getReturnValue() *Value {
	return t.returnValue
}

// ---------------- Return Statement ----------------
func (t *ReturnStmt) Execute(e *Environment) *Error {
	if executor := e.getReturnableTarget(); executor != nil {
		value, err := t.expr.Evaluate(e)
		if err != nil {
			return err
		}
		executor.setReturnValue(value)
		return nil
	}
	return NewError(t.token, "RuntimeError: 'continue' statement can only be used within an enclosing iteration")
}
