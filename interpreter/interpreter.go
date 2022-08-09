package interpreter

import (
	"fmt"
	"golox/parser"
	"golox/tokens"
)

func Interpret(expr parser.Expr) error {
    var i interpreter
    res := expr.Accept(&i)
    err, isError := res.(RuntimeException)
    if isError {
        return err
    }
    fmt.Println(res)
    return nil
}

type interpreter struct {

}

func (i *interpreter) VisitLiteral(l parser.Literal) interface{} {
    return l.Value
}

func (i *interpreter) VisitGrouping(g parser.Grouping) interface{} {
    if res, isError := g.Expression.Accept(i).(RuntimeException); isError {
        return res.Add("at grouping: ")
    } else {
        return res
    }
}

func (i *interpreter) VisitUnary(u parser.Unary) interface{} {
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

func (i *interpreter) VisitBinary(b parser.Binary) interface{} {
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
        l, lok := left.(float64)
        r, rok := right.(float64)
        if !lok || !rok {
            return NewRuntimeException(
                fmt.Sprintf("cannot preform '/' on %T and  %T", left, right))
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

