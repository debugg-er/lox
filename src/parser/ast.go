package parser

import (
	l "github.com/debugg-er/lox/src/lexer"
)

type (
	Expr interface {
		Expr()
	}

	Stmt interface {
		Stmt()
	}

	Callable interface {
		Call()
	}

	Loopable interface {
		SetBreaked(isBreaked bool)
		IsBreaked() bool
		SetContinued(isContinued bool)
		IsContinued() bool
	}

	Returnable interface {
		IsReturned() bool
		SetIsReturned(isReturned bool)
		// setReturnValue(value *Value)
		// getReturnValue() *Value
	}
)

type (
	PrimaryExpr struct {
		Value *l.Token
	}

	UnaryExpr struct {
		Operator *l.Token
		Operand  Expr
	}

	BinaryExpr struct {
		Operator *l.Token
		Left     Expr
		Right    Expr
	}

	VariableExpr struct {
		Name *l.Token
	}

	AssignExpr struct {
		Name  *l.Token
		Value Expr
	}

	FuncExpr struct {
		FuncStmt *FuncStmt
	}

	CallExpr struct {
		Callee    Expr
		Arguments []Expr
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
		Name       *l.Token
		Initilizer Expr
	}

	BlockStmt struct {
		Declarations []Stmt
	}

	IfStmt struct {
		Condition Expr
		ThenStmt  Stmt
		ElseStmt  Stmt
	}

	WhileStmt struct {
		Condition    Expr
		Body         Stmt
		_isBreaked   bool
		_isContinued bool
	}

	ForStmt struct {
		Initialization Stmt
		Condition      Expr
		Updation       Expr
		Body           Stmt
		_isBreaked     bool
		_isContinued   bool
	}

	BreakStmt struct {
		Token *l.Token
	}

	ContinueStmt struct {
		Token *l.Token
	}

	FuncStmt struct {
		Name        *l.Token
		Parameters  []*l.Token
		Body        *BlockStmt
		_isReturned bool
	}

	ReturnStmt struct {
		Token *l.Token
		Expr  Expr
	}
)

func (t *PrintStmt) Stmt()    {}
func (t *ExprStmt) Stmt()     {}
func (t *VarStmt) Stmt()      {}
func (t *BlockStmt) Stmt()    {}
func (t *IfStmt) Stmt()       {}
func (t *BreakStmt) Stmt()    {}
func (t *ContinueStmt) Stmt() {}
func (t *ReturnStmt) Stmt()   {}

func (t *WhileStmt) Stmt()                         {}
func (t *WhileStmt) IsBreaked() bool               { return t._isBreaked }
func (t *WhileStmt) SetBreaked(isBreaked bool)     { t._isBreaked = isBreaked }
func (t *WhileStmt) IsContinued() bool             { return t._isContinued }
func (t *WhileStmt) SetContinued(isContinued bool) { t._isContinued = isContinued }

func (t *ForStmt) Stmt()                         {}
func (t *ForStmt) IsBreaked() bool               { return t._isBreaked }
func (t *ForStmt) SetBreaked(isBreaked bool)     { t._isBreaked = isBreaked }
func (t *ForStmt) IsContinued() bool             { return t._isContinued }
func (t *ForStmt) SetContinued(isContinued bool) { t._isContinued = isContinued }

func (t *FuncStmt) Stmt()                         {}
func (t *FuncStmt) IsReturned() bool              { return t._isReturned }
func (t *FuncStmt) SetIsReturned(isReturned bool) { t._isReturned = isReturned }

// func (t *FuncStmt) setReturnValue(value *Value) {
// 	t.returnValue = value
// }

// func (t *FuncStmt) getReturnValue() *Value {
// 	return t.returnValue
// }

func (e *PrimaryExpr) Expr()  {}
func (e *UnaryExpr) Expr()    {}
func (e *BinaryExpr) Expr()   {}
func (e *VariableExpr) Expr() {}
func (e *AssignExpr) Expr()   {}
func (e *FuncExpr) Expr()     {}
func (e *CallExpr) Expr()     {}
