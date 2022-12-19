package lexer

import (
	"fmt"
	"strconv"
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

func NewLexer(source string) *Lexer {
	return &Lexer{
		current: 0,
		line:    0,
		source:  source,
		tokens:  make([]Token, 0),
	}
}

func (lexer *Lexer) Parse() []Token {
	for !lexer.isAtEnd() {
		lexer.scanToken()
	}
	return lexer.tokens
}

func (lexer *Lexer) scanToken() {
	c := lexer.advance()
	fmt.Println(c)

	switch c {
	case "(":
		lexer.addToken(LEFT_PAREN, nil)
		break
	case ")":
		lexer.addToken(RIGHT_PAREN, nil)
		break
	case "{":
		lexer.addToken(LEFT_BRACE, nil)
		break
	case "}":
		lexer.addToken(RIGHT_BRACE, nil)
		break
	case ",":
		lexer.addToken(COMMA, nil)
		break
	case ".":
		lexer.addToken(DOT, nil)
		break
	case "-":
		lexer.addToken(MINUS, nil)
		break
	case "+":
		lexer.addToken(PLUS, nil)
		break
	case "*":
		lexer.addToken(STAR, nil)
		break
	case ";":
		lexer.addToken(SEMICOLON, nil)
		break
	case "\r":
	case " ":
		break
	case "\n":
		lexer.line = lexer.line + 1
		break
	case "/":
		if lexer.match("/") {
			for lexer.advance() != "\n" {
			}
			lexer.line = lexer.line + 1
		} else {
			lexer.addToken(SLASH, nil)
		}
		break
	case ">":
		if lexer.match("=") {
			lexer.addToken(GREATER_EQUAL, nil)
		} else {
			lexer.addToken(GREATER, nil)
		}
		break
	case "<":
		if lexer.match("=") {
			lexer.addToken(LESS_EQUAL, nil)
		} else {
			lexer.addToken(LESS, nil)
		}
		break
	case "=":
		if lexer.match("=") {
			lexer.addToken(EQUAL_EQUAL, nil)
		} else {
			lexer.addToken(EQUAL, nil)
		}
		break
	case "!":
		if lexer.match("=") {
			lexer.addToken(BANG_EQUAL, nil)
		} else {
			lexer.addToken(BANG, nil)
		}
		break

	case "\"":
		lexer.string()
		break
	default:
		if isDigit(c) {
			lexer.number()
		}
		if isAlphabet(c) {
			lexer.word()
		}
		break
	}
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
	}

	value := lexer.source[start : lexer.current-1]
	lexer.addToken(STRING, value)
}

func (lexer *Lexer) number() {
	start := lexer.current - 1
	for isDigit(lexer.peek()) && !lexer.isAtEnd() {
		lexer.advance()
	}

	value := lexer.source[start : lexer.current-1]
	num, _ := strconv.Atoi(value)
	lexer.addToken(NUMBER, num)
}

func (lexer *Lexer) word() {
	word := string(lexer.source[lexer.current-1])
	for isAlphabet(lexer.peek()) && !lexer.isAtEnd() {
		word = word + lexer.advance()
		if keywords[word] != Undefined {
			lexer.addToken(keywords[word], nil)
			return
		}
	}

	lexer.addToken(IDENTIFIER, word)
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
	if c > "0" && c < "9" {
		return true
	}
	return false
}

func isAlphabet(c string) bool {
	if c > "a" && c < "z" {
		return true
	}
	if c > "A" && c < "Z" {
		return true
	}
	if c == "_" {
		return true
	}

	return false
}
