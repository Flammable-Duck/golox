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
	VisitAssign(a Assign) interface{}
    VisitVariable(v Variable) interface{}
}

type Literal struct {
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
	Operator   tokens.TokenType
}

func (u Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnary(u)
}

type Binary struct {
	Left     Expr
	Operator tokens.TokenType
	Right    Expr
}

func (b Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinary(b)
}

type Assign struct {
    Name tokens.Token
    Value Expr
}

func (a Assign) Accept(v ExprVisitor) interface{} {
    return v.VisitAssign(a)
}

type Variable struct {
	Name tokens.Token
}
func (v Variable) Accept(vis ExprVisitor) interface{} {
    return vis.VisitVariable(v)
}
