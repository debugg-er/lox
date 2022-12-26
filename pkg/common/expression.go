package common

type ExprType int

type Expr struct {
	Type     string
	Operator string
	Left     *Expr
	Right    *Expr
	Literal  interface{}
}
