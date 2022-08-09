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
	Literal  interface{}
}

func (t TokenType) String() string {
	var lexmes = map[TokenType]string{
		LeftParen:    "(",
		RightParen:   ")",
		LeftBrace:    "{",
		RightBrace:   "}",
		Comma:        ",",
		Dot:          ".",
		Minus:        "-",
		Plus:         "+",
		Semicolon:    ";",
		Slash:        "/",
		Star:         "*",
		Bang:         "!",
		BangEqual:    "!=",
		Equal:        "=",
		EqualEqual:   "==",
		Greater:      ">",
		GreaterEqual: ">=",
		Less:         "<",
		LessEqual:    "<=",
		And:          "and",
		Class:        "class",
		Else:         "Else",
		False:        "false",
		Fun:          "fun",
		For:          "for",
		If:           "if",
		Nil:          "Nil",
		Or:           "or",
		Print:        "print",
		Return:       "Return",
		Super:        "super",
		This:         "this",
		True:         "true",
		Var:          "var",
		While:        "while",
		Eof:          "<EOF>",
		Identifier:   "<identifier>",
		String:       "<string>",
		Number:       "<number>",
	}
	return lexmes[t]
}

func NewToken(t TokenType, p Position) Token {
	return NewTokenLiteral(t, t.String(), p)
}

func NewTokenLiteral(t TokenType, lexme string, p Position) Token {
	// return Token{Position: pos, Type: t, Lexeme: lexme}
	return Token{Type: t, Lexeme: lexme, Position: p}
}
