package common

import (
	"fmt"
	"os"
)

type Error struct {
	token   *Token
	message string
}

func (e *Error) Blame() {
	fmt.Fprintf(os.Stderr, "Line %d: %s\n", e.token.Line, e.message)
}

func NewError(token *Token, message string) *Error {
	return &Error{token, message}
}
