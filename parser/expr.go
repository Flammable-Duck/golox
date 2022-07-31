package parser

import (
	"golox/tokens"
)

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitLiteral(l Literal) interface{}
	VisitGrouping(g Grouping) interface{}
	VisitUnary(u Unary) interface{}
	VisitBinary(b Binary) interface{}
}

type Literal struct {
	// value tokens.TokenType
	Value interface{}
}

func (l Literal) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteral(l)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGrouping(g)
}

type Unary struct {
	Expression Expr
	Operator   tokens.Token
}

func (u Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnary(u)
}

type Binary struct {
	Left     Expr
	Operator tokens.Token
	Right    Expr
}

func (b Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinary(b)
}
