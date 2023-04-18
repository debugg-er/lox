package parser

import (
	"fmt"

	"github.com/debugg-er/lox/src/lexer"
)

type Error struct {
	token   *lexer.Token
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ParserError: Line %d at '%s': %s\n", e.token.Line, e.token.Type, e.message)
}

func NewParserError(token *lexer.Token, message string) *Error {
	return &Error{token, message}
}
