package parser

import (
	"fmt"
	"golox/tokens"
)

type Parser struct {
	tokens  []tokens.Token
	current int
}

func New(tokens []tokens.Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() Expr {
    expr := p.expression()
    return expr }

// rules

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(tokens.BangEqual, tokens.EqualEqual) {
		expr = Binary{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.comparison(),
		}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(tokens.Greater, tokens.GreaterEqual,
		tokens.Less, tokens.LessEqual) {
		expr = Binary{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.term(),
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(tokens.Minus, tokens.Plus) {
		expr = Binary{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.factor(),
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(tokens.Slash, tokens.Star) {
		expr = Binary{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.unary(),
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(tokens.Bang, tokens.Minus) {
		return Unary{
			Operator:   p.previous(),
			Expression: p.unary(),
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {

	var expr Expr
	if p.match(tokens.False) {
		expr = Literal{Value: false}

	} else if p.match(tokens.True) {
		expr = Literal{Value: true}

	} else if p.match(tokens.Nil) {
		expr = Literal{Value: nil}

	} else if p.match(tokens.Number, tokens.String) {
		expr = Literal{Value: p.previous().Literal}

	} else if p.match(tokens.LeftParen) {
		expr = p.expression()
		_, err := p.consume(tokens.RightParen, "Expect ')' after expression")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
    
	return expr
}

// error handling

func (p *Parser) consume(tpe tokens.TokenType, message string) (tokens.Token, error) {
	if p.check(tpe) {
		return p.advance(), nil
	}
	return p.advance(), &parseError{Token: p.advance(), Reason: message}
}

func (p *Parser) synchronize() {
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

func (p *Parser) match(types ...tokens.TokenType) bool {
	for _, tokentype := range types {
		if p.check(tokentype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p Parser) check(tokentype tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokentype
}

func (p Parser) peek() tokens.Token {
	return p.tokens[p.current]
}

func (p Parser) isAtEnd() bool {
	return p.peek().Type == tokens.Eof
}

func (p *Parser) advance() tokens.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() tokens.Token {
	return p.tokens[p.current-1]

}
