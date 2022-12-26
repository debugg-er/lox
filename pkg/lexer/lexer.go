package lexer

import (
	"fmt"
	"strconv"

	. "github.com/debugg-er/lox/pkg/common"
)

type Lexer struct {
	source  string
	current int
	line    int
	tokens  []Token
}

var keywords = map[string]TokenType{
	"var":    VAR,
	"and":    AND,
	"or":     OR,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"nil":    NIL,
	"for":    FOR,
	"while":  WHILE,
	"fun":    FUN,
	"print":  PRINT,
	"return": RETURN,
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
	tokens := lexer.tokens
	// Reset Lexer, keep line number
	lexer.tokens = make([]Token, 0)
	lexer.current = 0
	lexer.source = ""

	return tokens, nil
}

func (lexer *Lexer) scanToken() error {
	c := lexer.advance()

	switch c {
	case "(":
		lexer.addToken(LEFT_PAREN, nil)
	case ")":
		lexer.addToken(RIGHT_PAREN, nil)
	case "{":
		lexer.addToken(LEFT_BRACE, nil)
	case "}":
		lexer.addToken(RIGHT_BRACE, nil)
	case ",":
		lexer.addToken(COMMA, nil)
	case ".":
		lexer.addToken(DOT, nil)
	case "-":
		lexer.addToken(MINUS, nil)
	case "+":
		lexer.addToken(PLUS, nil)
	case "*":
		lexer.addToken(STAR, nil)
	case ";":
		lexer.addToken(SEMICOLON, nil)
	case "\r":
	case " ":
		break
	case "\n":
		lexer.line = lexer.line + 1
	case "/":
		if lexer.match("/") {
			for lexer.peek() != "\n" {
				lexer.advance()
			}
		} else {
			lexer.addToken(SLASH, nil)
		}
	case ">":
		if lexer.match("=") {
			lexer.addToken(GREATER_EQUAL, nil)
		} else {
			lexer.addToken(GREATER, nil)
		}
	case "<":
		if lexer.match("=") {
			lexer.addToken(LESS_EQUAL, nil)
		} else {
			lexer.addToken(LESS, nil)
		}
	case "=":
		if lexer.match("=") {
			lexer.addToken(EQUAL_EQUAL, nil)
		} else {
			lexer.addToken(EQUAL, nil)
		}
	case "!":
		if lexer.match("=") {
			lexer.addToken(BANG_EQUAL, nil)
		} else {
			lexer.addToken(BANG, nil)
		}

	case "\"":
		lexer.string()
	default:
		if isDigit(c) {
			lexer.number()
		} else if isAlphabet(c) {
			lexer.identifier()
		} else {
			return fmt.Errorf("unexpected token at line %d", lexer.line)
		}
	}

	return nil
}

func (lexer *Lexer) advance() string {
	c := string(lexer.source[lexer.current])
	lexer.current = lexer.current + 1
	return c
}

func (lexer *Lexer) peek() string {
	return string(lexer.source[lexer.current])
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current == len(lexer.source)
}

func (lexer *Lexer) match(c string) bool {
	if lexer.isAtEnd() {
		return false
	}
	if lexer.peek() != c {
		return false
	}
	lexer.advance()
	return true
}

func (lexer *Lexer) string() {
	start := lexer.current
	for !lexer.match("\"") && !lexer.isAtEnd() {
		lexer.advance()
	}

	value := lexer.source[start : lexer.current-1]
	lexer.addToken(STRING, value)
}

func (lexer *Lexer) number() {
	start := lexer.current - 1
	for isDigit(lexer.peek()) && !lexer.isAtEnd() {
		lexer.advance()
	}

	value := lexer.source[start:lexer.current]
	num, _ := strconv.Atoi(value)
	lexer.addToken(NUMBER, num)
}

func (lexer *Lexer) identifier() {
	word := string(lexer.source[lexer.current-1])
	for isAlphabet(lexer.peek()) && !lexer.isAtEnd() {
		word = word + lexer.advance()
	}

	if keywords[word] != Undefined {
		lexer.addToken(keywords[word], nil)
	} else {
		lexer.addToken(IDENTIFIER, word)
	}
}

func (lexer *Lexer) addToken(_type TokenType, value any) {
	token := Token{
		Type:  TokenType(_type),
		Value: value,
		Line:  lexer.line,
	}
	lexer.tokens = append(lexer.tokens, token)
}

func isDigit(c string) bool {
	if c >= "0" && c <= "9" {
		return true
	}
	return false
}

func isAlphabet(c string) bool {
	if c >= "a" && c <= "z" {
		return true
	}
	if c >= "A" && c <= "Z" {
		return true
	}
	if c == "_" {
		return true
	}

	return false
}
