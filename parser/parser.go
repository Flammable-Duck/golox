package parser

import (
	"fmt"
	"golox/tokens"
)

type parser struct {
	tokens  []tokens.Token
	current int
}

func Parse(tokens []tokens.Token) []Stmt {
	p := parser{tokens: tokens, current: 0}
	var statments []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			fmt.Printf("\u001b[31m%s\u001b[39m\n", err.Error())
			p.synchronize()
		} else {
			statments = append(statments, stmt)
		}
	}
	return statments
}

// rules

func (p *parser) declaration() (Stmt, error) {
	if p.match(tokens.Var) {
		return p.varDeclaration()
	}
	return p.statment()
}

func (p *parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(tokens.Identifier, "expected variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(tokens.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after declaration.")
	if err != nil {
		return nil, err
	}
	return Var{Name: name, Initializer: initializer}, nil
}

func (p *parser) statment() (Stmt, error) {
	if p.match(tokens.Print) {
		return p.printStatement()
	}

	res, err := p.expressionStatement()
	return res, err
}

func (p *parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after expression.")
	if err != nil {
		return nil, err
	}

	return PrintStmt{Expression: value}, nil
}

func (p *parser) expressionStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after expression.")
	if err != nil {
		return nil, err
	}
	return ExprStmt{Expression: value}, nil
}

func (p *parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if p.match(tokens.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if expr, ok := expr.(Variable); ok {
			return Assign{Name: expr.Name, Value: value}, nil
		} else {
			return nil, parseError{
				Token:  equals,
				Reason: "Invialid Assignment target."}
		}
	}

	return expr, err
}

func (p *parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return expr, err
	}
	for p.match(tokens.BangEqual, tokens.EqualEqual) {
		op := p.previous().Type
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return expr, err
	}

	for p.match(tokens.Greater, tokens.GreaterEqual,
		tokens.Less, tokens.LessEqual) {
		op := p.previous().Type
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return expr, err
	}

	for p.match(tokens.Minus, tokens.Plus) {
		op := p.previous().Type
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return expr, err
	}

	for p.match(tokens.Slash, tokens.Star) {
		op := p.previous().Type
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (Expr, error) {
	if p.match(tokens.Bang, tokens.Minus) {
		op := p.previous().Type
		expr, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{
			Operator:   op,
			Expression: expr,
		}, nil
	}
	return p.primary()
}

func (p *parser) primary() (Expr, error) {

	var expr Expr
	if p.match(tokens.False) {
		expr = Literal{Value: false}

	} else if p.match(tokens.True) {
		expr = Literal{Value: true}

	} else if p.match(tokens.Nil) {
		expr = Literal{Value: nil}

	} else if p.match(tokens.Number, tokens.String) {
		expr = Literal{Value: p.previous().Literal}

	} else if p.match(tokens.Identifier) {
		expr = Variable{Name: p.previous()}

	} else if p.match(tokens.LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return expr, err
		}
		_, err = p.consume(tokens.RightParen, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

// error handling

func (p *parser) consume(tpe tokens.TokenType, message string) (tokens.Token, error) {
	if p.check(tpe) {
		return p.advance(), nil
	}
	return p.advance(), &parseError{Token: p.advance(), Reason: message}
}

func (p *parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == tokens.Semicolon {
			return
		}
		switch p.peek().Type {
		case tokens.Class, tokens.For, tokens.Fun, tokens.If, tokens.Print,
			tokens.Return, tokens.Var, tokens.While:
			return
		}
		p.advance()
	}
}

// utilty methods

func (p *parser) match(types ...tokens.TokenType) bool {
	for _, tokentype := range types {
		if p.check(tokentype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p parser) check(tokentype tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokentype
}

func (p parser) peek() tokens.Token {
	return p.tokens[p.current]
}

func (p parser) isAtEnd() bool {
	return p.peek().Type == tokens.Eof
}

func (p *parser) advance() tokens.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) previous() tokens.Token {
	return p.tokens[p.current-1]

}
