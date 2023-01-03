package parser

type StmtType int

const (
	PRINT_STMM = iota
	EXPR_STMT
	IF_STMT
	FOR_STMT
	WHILE_STMT
	RETURN_STMT
)
