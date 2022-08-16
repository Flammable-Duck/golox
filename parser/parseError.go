package parser

import (
	"fmt"
	"golox/tokens"
)

type parseError struct {
	Token  tokens.Token
	Reason string
}

func (e parseError) Error() string {
	if e.Token.Type == tokens.Eof {
		return fmt.Sprintf("%d:%d parse error at end: %s",
			e.Token.Position.Row,
			e.Token.Position.Col,
			e.Reason)

	} else {
		return fmt.Sprintf("%d:%d parse error near '%s': %s",
			e.Token.Position.Row,
			e.Token.Position.Col,
			e.Token.Lexeme,
			e.Reason)
	}
}
