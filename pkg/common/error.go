package common

import (
	"fmt"
)

type Error struct {
	token   *Token
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Line %d at '%s': %s\n", e.token.Line, e.token.Type, e.message)
}

func NewError(token *Token, message string) *Error {
	return &Error{token, message}
}
