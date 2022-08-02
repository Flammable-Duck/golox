package scanner

import (
	"fmt"
	"golox/tokens"
	"strconv"
	"strings"
)

type Scanner struct {
	src     []rune
	tokens  []tokens.Token
	errors  []error
	start   int
	current int
    scanError struct {}
}


func NewScanner(source string) Scanner {
	return Scanner{
		src:     []rune(source),
		start:   0,
		current: 0,
	}
}

func (s *Scanner) Read() (tokens.Token, error) {
	var tok tokens.Token
	var tokenFound bool = false
	var err error = nil
	if s.isAtEnd() {
		return s.newToken(tokens.Eof), err
	}

	s.start = s.current
	r := s.advance()

	switch r {
	case '(':
		tok = s.newToken(tokens.LeftParen)
		tokenFound = true
	case ')':
		tok = s.newToken(tokens.RightParen)
		tokenFound = true
	case '{':
		tok = s.newToken(tokens.LeftBrace)
		tokenFound = true
	case '}':
		tok = s.newToken(tokens.RightBrace)
		tokenFound = true
	case ',':
		tok = s.newToken(tokens.Comma)
		tokenFound = true
	case '.':
		tok = s.newToken(tokens.Dot)
		tokenFound = true
	case '-':
		tok = s.newToken(tokens.Minus)
		tokenFound = true
	case '+':
		tok = s.newToken(tokens.Plus)
		tokenFound = true
	case ';':
		tok = s.newToken(tokens.Semicolon)
		tokenFound = true
	case '*':
		tok = s.newToken(tokens.Star)
		tokenFound = true
	case '!':
		if s.match('=') {
			tok = s.newToken(tokens.BangEqual)
			tokenFound = true
		} else {
			tok = s.newToken(tokens.Bang)
			tokenFound = true
		}
	case '=':
		if s.match('=') {
			tok = s.newToken(tokens.EqualEqual)
			tokenFound = true
		} else {
			tok = s.newToken(tokens.Equal)
			tokenFound = true
		}
	case '<':
		if s.match('=') {
			tok = s.newToken(tokens.LessEqual)
			tokenFound = true
		} else {
			tok = s.newToken(tokens.Less)
			tokenFound = true
		}
	case '>':
		if s.match('=') {
			tok = s.newToken(tokens.GreaterEqual)
			tokenFound = true
		} else {
			tok = s.newToken(tokens.Greater)
			tokenFound = true
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			l := 1
			for l != 0 {
				if s.isAtEnd() {
					// err = fmt.Errorf("unclosed comment")
					tok = s.newToken(tokens.Eof)
					tokenFound = true
					err = s.error("unclosed comment")
					break
				}
				c := s.advance()
				switch c {
				case '/':
					if s.match('*') {
						l++
					}
				case '*':
					if s.match('/') {
						l--
					}
				}
			}
		} else {
			tok = s.newToken(tokens.Slash)
			tokenFound = true
		}
	case '"':
		tok, err = s.string()
		tokenFound = true
	case ' ', '\t':
		// tok, err = s.Read()
	case '\n':
		// tok, err = s.Read()
	default:
		if isDigit(r) {
			tok, err = s.number()
			tokenFound = true
		} else if isAlpha(r) {
			tok, err = s.identifier()
			tokenFound = true
		} else {
			err = s.error(fmt.Sprintf("Unrecognized character: %s", string(r)))
		}
	}

	// if we were not able to find any tokens, make a second pass
	if !tokenFound {
		if err == nil {
			tok, err = s.Read()
		} else {
			tok, _ = s.Read()
		}
	}

	return tok, err
}

func (s *Scanner) string() (tokens.Token, error) {
	var err error
	var value string = ""
	for s.peek() != '"' && !s.isAtEnd() {
		s.advance()
	}
	if s.isAtEnd() && s.peek() != '"' {
		err = s.error("unterminated string")
        return tokens.Token{}, err
	} else {
        s.advance()
    }
	value = string(s.src[s.start+1 : s.current-1])
	tok := tokens.Token{
		Position: s.getPosition(),
		Type:     tokens.String,
		Lexeme:   value,
		Literal:  value,
	}
	return tok, err
}

func (s *Scanner) number() (tokens.Token, error) {
	var err error = nil
	for isDigit(s.peek()) && !s.isAtEnd() {
		s.advance()
	}
	if s.peek() == '.' {
		s.advance()
	}
	for isDigit(s.peek()) && !s.isAtEnd() {
		s.advance()
	}
	// s.advance()
    lexme := string(s.src[s.start:s.current])
    var literal interface{}
    if strings.ContainsRune(lexme, '.') {
        literal, _ = strconv.ParseFloat(lexme, 64)
    } else {
        literal, _ = strconv.Atoi(lexme)
    }

	tok := tokens.Token{
		Position: s.getPosition(),
		Type:     tokens.String,
		Lexeme:   lexme,
		Literal:  literal,
	}
	return tok, err
}

func (s *Scanner) identifier() (tokens.Token, error) {
	var err error = nil
	var value string = ""
	var tok tokens.Token
	var identifiers = map[string]tokens.TokenType{
		"and":    tokens.And,
		"class":  tokens.Class,
		"else":   tokens.Else,
		"false":  tokens.False,
		"for":    tokens.For,
		"fun":    tokens.Fun,
		"if":     tokens.If,
		"nil":    tokens.Nil,
		"or":     tokens.Or,
		"print":  tokens.Print,
		"return": tokens.Return,
		"super":  tokens.Super,
		"this":   tokens.This,
		"true":   tokens.True,
		"while":  tokens.While,
	}

	for isAlpha(s.peek()) && !s.isAtEnd() {
		s.advance()
	}

	value = string(s.src[s.start:s.current])
	keyword, exists := identifiers[value]
	if exists {
		tok = tokens.NewTokenLiteral(keyword, value,
			s.getPosition())
	} else {
		tok = tokens.NewTokenLiteral(tokens.Identifier, value,
			s.getPosition())
	}
	return tok, err
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\u0000'
	}
	return s.src[s.current]
}

func (s *Scanner) match(expected rune) bool {
	// fmt.Printf("Is the next token \"%s\"? ", string(expected))
	if s.isAtEnd() {
		// fmt.Println("Dont bother checking, were at the end.")
		return false
	}

	r := s.src[s.current]
	if r != expected {
		// fmt.Println("Naw.")
		return false
	}
	// fmt.Println("Yeah!")
	s.current++
	return true
}

func (s *Scanner) advance() rune {
	r := s.src[s.current]
	s.current++
	return r
}

func (s *Scanner) error(e string) error {
	pos := s.getPosition()
	err := fmt.Errorf("%d:%d: %s", pos.Row, pos.Col, e)
	s.errors = append(s.errors, err)
	return err
}

// wrapper for tokens.newToken
func (s *Scanner) newToken(t tokens.TokenType) tokens.Token {
	return tokens.NewToken(t, s.getPosition())
}

// this is a slow and dumb and bad way to do this but ¯\_(ツ)_/¯
func (s *Scanner) getPosition() tokens.Position {
	p := tokens.Position{Row: 1, Col: 1}
	for i := 0; i < s.start; i++ {
		if s.src[i] == '\n' {
			p.Row++
			p.Col = 1
		} else {
			p.Col++
		}
	}
	return p
}

func (s *Scanner) Errors() []error {
	return s.errors
}

func isAlpha(r rune) bool {
	if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' {
		return true
	}
	return false
}

func isDigit(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	return false
}
func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
