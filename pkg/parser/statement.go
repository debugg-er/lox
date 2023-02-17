package parser

import (
	"fmt"

	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type (
	Stmt interface {
		Execute(e *Environment) error
	}

	Loopable interface {
		setBreaked(isBreaked bool)
		isBreaked() bool
		setContinued(isContinued bool)
		isContinued() bool
	}

	Returnable interface {
		isReturned() bool
		setIsReturned(isReturned bool)
		setReturnValue(value *Value)
		getReturnValue() *Value
	}
)

type (
	PrintStmt struct {
		Expr Expr
	}

	ExprStmt struct {
		Expr Expr
	}

	VarStmt struct {
		name       *Token
		initilizer Expr
	}

	BlockStmt struct {
		declarations []Stmt
	}

	IfStmt struct {
		condition Expr
		thenStmt  Stmt
		elseStmt  Stmt
	}

	WhileStmt struct {
		condition    Expr
		body         Stmt
		_isBreaked   bool
		_isContinued bool
	}

	ForStmt struct {
		initialization Stmt
		condition      Expr
		updation       Expr
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
		parameters  []*Token
		body        *BlockStmt
		returnValue *Value
		_isReturned bool
	}

	ReturnStmt struct {
		token *Token
		expr  Expr
	}
)

// ---------------- Print Statement ----------------
func (t *PrintStmt) Execute(e *Environment) error {
	value, err := t.Expr.Evaluate(e)
	if err != nil {
		return err
	}
	fmt.Println(value.Stringify())
	return nil
}

// ---------------- Expression Statement ----------------
func (t *ExprStmt) Execute(e *Environment) error {
	_, err := t.Expr.Evaluate(e)
	return err
}

// ---------------- Variable Declaration Statement ----------------
func (t *VarStmt) Execute(e *Environment) error {
	value, err := t.initilizer.Evaluate(e)
	if err != nil {
		return err
	}
	e.define(t.name, value)
	return nil
}

// ---------------- Block Statement ----------------
func (t *BlockStmt) Execute(e *Environment) error {
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
		// Stop execute when return is met on child statements
		if executor := e.getReturnableTarget(); executor != nil {
			if executor.isReturned() {
				return nil
			}
		}
	}
	return nil
}

// ---------------- If Statement ----------------
func (t *IfStmt) Execute(e *Environment) error {
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
func (t *WhileStmt) Execute(e *Environment) error {
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
		if t.isBreaked() {
			return nil
		}
		if t.isContinued() {
			t.setContinued(false)
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
func (t *ForStmt) Execute(e *Environment) error {
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
		if t.isBreaked() {
			return nil
		}
		if t.isContinued() {
			t.setContinued(false)
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
func (t *BreakStmt) Execute(e *Environment) error {
	executor := e.getLoopableTarget()
	if executor != nil {
		return NewError(t.token, "RuntimeError: 'break' statement can only be used within an enclosing iteration")
	}
	executor.setBreaked(true)
	return nil
}

// ---------------- Continue Statement ----------------
func (t *ContinueStmt) Execute(e *Environment) error {
	executor := e.getLoopableTarget()
	if executor == nil {
		return NewError(t.token, "RuntimeError: 'continue' statement can only be used within an enclosing iteration")
	}
	executor.setContinued(true)
	return nil
}

// ---------------- Function Statement ----------------
func (t *FuncStmt) Execute(e *Environment) error {
	e.returableTarget = t
	err := t.body.Execute(e)
	if err != nil {
		return err
	}
	return nil
}

func (t *FuncStmt) setReturnValue(value *Value) {
	t.returnValue = value
}

func (t *FuncStmt) getReturnValue() *Value {
	return t.returnValue
}

func (t *FuncStmt) isReturned() bool {
	return t._isReturned
}

func (t *FuncStmt) setIsReturned(isReturned bool) {
	t._isReturned = isReturned
}

// ---------------- Return Statement ----------------
func (t *ReturnStmt) Execute(e *Environment) error {
	executor := e.getReturnableTarget()
	if executor == nil {
		return NewError(t.token, "RuntimeError: 'continue' statement can only be used within an enclosing iteration")
	}
	value, err := t.expr.Evaluate(e)
	if err != nil {
		return err
	}
	executor.setIsReturned(true)
	executor.setReturnValue(value)
	return nil
}
