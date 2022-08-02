package tokens

type TokenType int

const (
	// single character tokens
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// one or two character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// literals
	Identifier
	String
	Number

	// keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	Eof
)

type Position struct {
	Row int
	Col int
}

type Token struct {
	Position Position
	Type     TokenType
	Lexeme   string
    Literal interface{}
}

func stringFromTokenType(t TokenType) string {
	var lexme string
	switch t {
	case LeftParen:
		lexme = "("
	case RightParen:
		lexme = ")"
	case LeftBrace:
		lexme = "{"
	case RightBrace:
		lexme = "}"
	case Comma:
		lexme = ","
	case Dot:
		lexme = "."
	case Minus:
		lexme = "-"
	case Plus:
		lexme = "+"
	case Semicolon:
		lexme = ";"
	case Slash:
		lexme = "/"
	case Star:
		lexme = "*"
	case Bang:
		lexme = "!"
	case BangEqual:
		lexme = "!="
	case Equal:
		lexme = "="
	case EqualEqual:
		lexme = "=="
	case Greater:
		lexme = ">"
	case GreaterEqual:
		lexme = ">="
	case Less:
		lexme = "<"
	case LessEqual:
		lexme = "<="
	case And:
		lexme = "and"
	case Class:
		lexme = "class"
	case Else:
		lexme = "Else"
	case False:
		lexme = "false"
	case Fun:
		lexme = "fun"
	case For:
		lexme = "for"
	case If:
		lexme = "if"
	case Nil:
		lexme = "Nil"
	case Or:
		lexme = "or"
	case Print:
		lexme = "print"
	case Return:
		lexme = "Return"
	case Super:
		lexme = "super"
	case This:
		lexme = "this"
	case True:
		lexme = "true"
	case Var:
		lexme = "var"
	case While:
		lexme = "while"
	case Eof:
		lexme = "\u0000"
	case Identifier:
        lexme= "<identifier>"
    case String:
        lexme= "<string>"
    case Number:
        lexme= "<number>"
	}
    return lexme
}

func NewToken(t TokenType, p Position) Token {
	return NewTokenLiteral(t, stringFromTokenType(t), p)
}

func NewTokenLiteral(t TokenType, lexme string, p Position) Token {
	// return Token{Position: pos, Type: t, Lexeme: lexme}
	return Token{Type: t, Lexeme: lexme, Position: p}
}
