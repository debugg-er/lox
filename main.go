package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/debugg-er/lox/pkg/lexer"
)

func main() {
	if len(os.Args) > 1 {
		execFile()
	} else {
		enterPrompt()
	}
}

func execFile() {
	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "File not found")
		os.Exit(1)
	}
	var scanner Test = lexer.NewLexer()
	execute(string(source), scanner)
}

func enterPrompt() {
	reader := bufio.NewReader(os.Stdin)
	scanner := lexer.NewLexer()
	for {
		fmt.Print("> ")
		code, _ := reader.ReadString('\n')
		execute(code, scanner)
	}
}

func execute(source string, scanner Test) {
	tokens, err := scanner.Parse(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
}

type Test interface {
	Parse(source string) ([]lexer.Token, error)
}
