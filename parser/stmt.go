package parser

import (
	"golox/tokens"
)

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitPrintStmt(prnt PrintStmt) interface{}
	VisitExprStmt(expr ExprStmt) interface{}
	VisitVarStmt(Var) interface{}
}

type ExprStmt struct {
	Expression Expr
}

func (e ExprStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitExprStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (p PrintStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(p)
}

type Var struct {
	Name        tokens.Token
	Initializer Expr
}

func (v Var) Accept(vis StmtVisitor) interface{} {
	return vis.VisitVarStmt(v)
}
