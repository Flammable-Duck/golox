package parser

type Stmt interface {
    Accept(v StmtVisitor) interface{}
}

type StmtVisitor interface {
    VisitPrintStmt(prnt PrintStmt)    interface{}
    VisitExprStmt(expr ExprStmt) interface{}

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
