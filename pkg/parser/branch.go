package parser

import (
	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type context struct {
	inFor      bool
	inWhile    bool
	inFunction bool
}

func verifyBranching(stmt Stmt) []*Error {
	context := &context{}

	switch stmt.(type) {
	case *BlockStmt, *IfStmt, *ForStmt, *WhileStmt:
		return _verifyBranching(stmt, context)
	default:
		return nil
	}
}

func _verifyBranching(stmt Stmt, context *context) []*Error {
	switch stmt := stmt.(type) {
	case *BreakStmt:
		if !context.inFor && !context.inWhile {
			return []*Error{NewError(stmt.token, "SynctaxError: 'break' statement can only be used within an enclosing iteration")}
		}
	case *ContinueStmt:
		if !context.inFor && !context.inWhile {
			return []*Error{NewError(stmt.token, "SynctaxError: 'continue' statement can only be used within an enclosing iteration")}
		}
	case *ReturnStmt:
		if !context.inFunction {
			return []*Error{NewError(stmt.token, "SynctaxError: 'return' statement can only be used within function")}
		}
	case *ForStmt:
		context.inFor = true
		return _verifyBranching(stmt.body, context)
	case *WhileStmt:
		context.inWhile = true
		return _verifyBranching(stmt.body, context)
	case *FuncStmt:
		context.inFunction = true
		return _verifyBranching(stmt.body, context)
	case *IfStmt:
		return _verifyBranching(stmt.thenStmt, context)
	case *BlockStmt:
		errors := make([]*Error, 0)
		for _, childStmt := range stmt.declarations {
			errors = append(errors, _verifyBranching(childStmt, context)...)
		}
		if len(errors) == 0 {
			return nil
		}
		return errors
	}

	return nil
}
