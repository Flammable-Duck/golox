package astPrinter

import (
	"fmt"
	"golox/parser"
)

type astPrinter struct {
    depth int
}

func (p *astPrinter) VisitLiteral(l parser.Literal) interface{} {
    switch l.Value.(type) {
    case nil:
        return "nil"
    case string:
        return fmt.Sprintf("\"%s\"", l.Value)
    default:
        return fmt.Sprintf("%v", l.Value)
    }
}

func (p *astPrinter) VisitGrouping(g parser.Grouping) interface{} {
    str := fmt.Sprintf("(\n%s)\n", g.Expression.Accept(p))
    p.depth++
    return str
}

func (p *astPrinter) VisitUnary(u parser.Unary) interface{} {
    str := fmt.Sprintf("(%s %s)\n",
        u.Operator.String(),
        u.Expression.Accept(p))
    p.depth--
    return str
}

func (p *astPrinter) VisitBinary(b parser.Binary) interface{} {
    p.depth++
    str := fmt.Sprintf("(%s %s",
        b.Operator.String(),
        b.Left.Accept(p))

    switch b.Right.(type) {
        case parser.Literal:
            str = fmt.Sprintf("%s %s)", str, b.Right.Accept(p))
        default:
            str = fmt.Sprintf("%s \n%s%s", str,
            indent(p.depth),
            b.Right.Accept(p))
            p.depth--
            str = fmt.Sprintf("%s\n%s)", str, indent(p.depth))
    }
    p.depth--
    return str
}

func (p *astPrinter) VisitPrintStmt(psr parser.PrintStmt) interface{} {
    return fmt.Sprintf("print %s\n", psr.Expression.Accept(p))
}

func (p *astPrinter) VisitExprStmt(exprStmt parser.ExprStmt) interface{} {
    return exprStmt.Expression.Accept(p)
}



func PrintAst(expr parser.Expr)  {
    fmt.Printf("%s\n", expr.Accept(&astPrinter{}))
}

func PrintStmt(stmt parser.Stmt) {
    fmt.Printf("%s\n", stmt.Accept(&astPrinter{}))
}

func indent(depth int) string{
    indent := ""
    for i := 0; i <= depth; i++ {
        indent = indent + "  "
    }
    return indent
}
