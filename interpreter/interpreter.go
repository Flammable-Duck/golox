package interpreter

import (
	"fmt"
	"golox/interpreter/environment"
	"golox/parser"
	"golox/tokens"
)

type Interpreter struct {
    env environment.Environment
}

func New() Interpreter {
    return Interpreter{env: environment.New()}
}

func (i Interpreter) Interpret(stmt parser.Stmt)  (interface{}, error) {
    res := stmt.Accept(i)
    err, isError := res.(RuntimeException)
    if isError {
        return nil, err
    }
    // fmt.Printf("\u001b[2m] %v\u001b[0m\n", res)
    return res, nil
}

func (i Interpreter) VisitPrintStmt(prnt parser.PrintStmt) interface{} {
    value := prnt.Expression.Accept(i)
    fmt.Printf("%v\n", value)
    return nil
}

func (i Interpreter) VisitExprStmt(exprstmt parser.ExprStmt) interface{} {
    return exprstmt.Expression.Accept(i)
}

func (i Interpreter) VisitVarStmt(vr parser.Var) interface{} {
    initalizer := vr.Initializer.Accept(i)
    i.env.Define(vr.Name.Lexeme, initalizer)
    return initalizer
}

func (i Interpreter) VisitVariable(v parser.Variable) interface{} {
    val, err := i.env.Get(v.Name.Lexeme)
    if err != nil {
        return NewRuntimeException(err.Error())
    }

    return val
}



func (i Interpreter) VisitLiteral(l parser.Literal) interface{} {
    return l.Value
}

func (i Interpreter) VisitGrouping(g parser.Grouping) interface{} {
    if res, isError := g.Expression.Accept(i).(RuntimeException); isError {
        return res.Add("at grouping: ")
    } else {
        return res
    }
}

func (i Interpreter) VisitUnary(u parser.Unary) interface{} {
    right := u.Expression.Accept(i)
    if res, isError := right.(RuntimeException); isError {
        return res.Add(fmt.Sprintf("at %s: ", u.Operator.String()))
    }
    switch u.Operator {
    case tokens.Minus:
            v, ok := right.(float64)
            if !ok {
                return NewRuntimeException("cannot perform '-' on string")
            }
            return -v
    case tokens.Bang:
            return !isTruthy(right)
    }
    return NewRuntimeException(fmt.Sprintf("unexpected operator: %s",
        u.Operator.String()))
}

func (i Interpreter) VisitBinary(b parser.Binary) interface{} {
    left := b.Left.Accept(i)
    if res, isError := left.(RuntimeException); isError {
        return res.Add(fmt.Sprintf("at %s: ", b.Operator.String()))
    }
    right := b.Right.Accept(i)
    if res, isError := right.(RuntimeException); isError {
        return res.Add(fmt.Sprintf("at %s: ", b.Operator.String()))
    }

    switch b.Operator {
    case tokens.Minus:
        l, lok := left.(float64)
        r, rok := right.(float64)
        if !lok || !rok {
            return NewRuntimeException("cannot preform '-' on string")
        }
        return l - r

    case tokens.Slash:
        l, lOk := left.(float64)
        r, rOk := right.(float64)
        if !lOk || !rOk {
            return NewRuntimeException(
                fmt.Sprintf("cannot preform '/' on %T and  %T", left, right))
        }
        if r == 0 {
            return NewRuntimeException("cannot divide by zero")
        }
        return l / r

    case tokens.Star:
        l, lok := left.(float64)
        r, rok := right.(float64)
        if !lok || !rok {
            return NewRuntimeException(
                fmt.Sprintf("cannot preform '/' on %T and  %T", left, right))
        }
        return l * r

    case tokens.Plus:
        switch l := left.(type) {
        case float64:
            r, ok := right.(float64)
            if !ok {
                return NewRuntimeException(
                    fmt.Sprintf("cannot preform '+' on %T and  %T", left, right))
            }
            return l + r
        case string:
            r, ok := right.(string)
            if !ok {
                return NewRuntimeException(
                    fmt.Sprintf("cannot preform '+' on %T and  %T", left, right))
            }
            return l + r
        }
    case tokens.Greater:
        l, lok := left.(float64)
        r, rok := right.(float64)
        if !lok || !rok {
            return NewRuntimeException(
                fmt.Sprintf("cannot preform '>' on %T and  %T", left, right))
        }
        return l > r
    case tokens.GreaterEqual:
        l, lok := left.(float64)
        r, rok := right.(float64)
        if !lok || !rok {
            return NewRuntimeException(
                fmt.Sprintf("cannot preform '>=' on %T and  %T", left, right))
        }
        return l >= r
    case tokens.Less:
        l, lok := left.(float64)
        r, rok := right.(float64)
        if !lok || !rok {
            return NewRuntimeException(
                fmt.Sprintf("cannot preform '<' on %T and  %T", left, right))
        }
        return l < r
    case tokens.LessEqual:
        l, lok := left.(float64)
        r, rok := right.(float64)
        if !lok || !rok {
            return NewRuntimeException(
                fmt.Sprintf("cannot preform '<=' on %T and  %T", left, right))
        }
        return l <= r
    case tokens.EqualEqual:
        return left == right
    case tokens.BangEqual:
        return left != right
    }
    return NewRuntimeException(
        fmt.Sprintf("unexpected operator: %s", b.Operator.String()))
}

func isTruthy(val any) bool {
    switch t := val.(type) {
    case bool:
        return t
    case nil:
        return false
    default:
        return true
    }
}

