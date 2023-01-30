package lexer

import (
	"bytes"
	"fmt"
	"strconv"

	//lint:ignore ST1001 that's what we want
	. "github.com/debugg-er/lox/pkg/common"
)

type Lexer struct {
	source  string
	current int
	line    int
	tokens  []Token
}

func NewLexer() *Lexer {
	return &Lexer{
		current: 0,
		line:    1,
		source:  "",
		tokens:  make([]Token, 0),
	}
}

func (lexer *Lexer) Parse(source string) ([]Token, error) {
	lexer.source = source
	for !lexer.isAtEnd() {
		err := lexer.scanToken()
		if err != nil {
			return nil, err
		}
	}
	lexer.addToken(EOF, nil)
	return lexer.tokens, nil
}

func (lexer *Lexer) scanToken() error {
	c := lexer.advance()

	switch c {
	case '(':
		lexer.addToken(LEFT_PAREN, nil)
	case ')':
		lexer.addToken(RIGHT_PAREN, nil)
	case '{':
		lexer.addToken(LEFT_BRACE, nil)
	case '}':
		lexer.addToken(RIGHT_BRACE, nil)
	case ',':
		lexer.addToken(COMMA, nil)
	case '.':
		lexer.addToken(DOT, nil)
	case '-':
		lexer.addToken(MINUS, nil)
	case '+':
		lexer.addToken(PLUS, nil)
	case '*':
		lexer.addToken(STAR, nil)
	case ';':
		lexer.addToken(SEMICOLON, nil)
	case '\r':
	case ' ':
		break
	case '\n':
		lexer.line = lexer.line + 1
	case '/':
		if lexer.match('/') {
			for !lexer.isAtEnd() && lexer.peek() != '\n' {
				lexer.advance()
			}
		} else {
			lexer.addToken(SLASH, nil)
		}
	case '>':
		if lexer.match('=') {
			lexer.addToken(GREATER_EQUAL, nil)
		} else {
			lexer.addToken(GREATER, nil)
		}
	case '<':
		if lexer.match('=') {
			lexer.addToken(LESS_EQUAL, nil)
		} else {
			lexer.addToken(LESS, nil)
		}
	case '=':
		if lexer.match('=') {
			lexer.addToken(EQUAL_EQUAL, nil)
		} else {
			lexer.addToken(EQUAL, nil)
		}
	case '!':
		if lexer.match('=') {
			lexer.addToken(BANG_EQUAL, nil)
		} else {
			lexer.addToken(BANG, nil)
		}

	case '"':
		err := lexer.string()
		if err != nil {
			return err
		}
	default:
		if isDigit(c) {
			err := lexer.number()
			if err != nil {
				return err
			}
		} else if isAlphabet(c) {
			lexer.identifier()
		} else {
			return fmt.Errorf("SyntaxError: Unexpected token at line %d", lexer.line)
		}
	}

	return nil
}

func (lexer *Lexer) string() error {
	var str bytes.Buffer
	for !lexer.isAtEnd() && !lexer.match('"') {
		c := lexer.advance()
		if c == '\\' {
			str.WriteByte(escapeSequence(lexer.advance()))
		} else {
			str.WriteByte(c)
		}
	}
	if lexer.isAtEnd() && lexer.previous() != '"' {
		return fmt.Errorf("SyntaxError: Expected '\"' at line %d", lexer.line)
	}

	lexer.addToken(STRING, str.String())
	return nil
}

func (lexer *Lexer) number() error {
	start := lexer.current - 1
	for !lexer.isAtEnd() && (isDigit(lexer.peek()) || lexer.peek() == '.') {
		lexer.advance()
	}

	value := lexer.source[start:lexer.current]
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("SyntaxError: Unexpected token at line %d", lexer.line)
	}
	lexer.addToken(NUMBER, num)
	return nil
}

func (lexer *Lexer) identifier() {
	start := lexer.current - 1
	for isAlphabet(lexer.peek()) && !lexer.isAtEnd() {
		lexer.advance()
	}

	identifier := lexer.source[start:lexer.current]
	if Keywords[identifier] != Undefined {
		var value any = nil
		if identifier == TRUE {
			value = true
		} else if identifier == FALSE {
			value = false
		}
		lexer.addToken(Keywords[identifier], value)
	} else {
		lexer.addToken(IDENTIFIER, identifier)
	}
}

func (lexer *Lexer) addToken(_type TokenType, value interface{}) {
	token := Token{
		Type:  TokenType(_type),
		Value: value,
		Line:  lexer.line,
	}
	lexer.tokens = append(lexer.tokens, token)
}

func (lexer *Lexer) advance() byte {
	c := lexer.source[lexer.current]
	lexer.current = lexer.current + 1
	return c
}

func (lexer *Lexer) peek() byte {
	return lexer.source[lexer.current]
}

func (lexer *Lexer) previous() byte {
	return lexer.source[lexer.current-1]
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current == len(lexer.source)
}

func (lexer *Lexer) match(c byte) bool {
	if lexer.isAtEnd() {
		return false
	}
	if lexer.peek() != c {
		return false
	}
	lexer.advance()
	return true
}

func isDigit(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func isAlphabet(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c == '_':
		return true
	default:
		return false
	}
}

func escapeSequence(c byte) byte {
	switch c {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	default:
		return c
	}
}
