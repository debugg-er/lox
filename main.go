package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/debugg-er/lox/pkg/lexer"
	"github.com/debugg-er/lox/pkg/parser"
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
	execute(string(source))
}

func enterPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		code, _ := reader.ReadString('\n')
		execute(code)
	}
}

func execute(source string) {
	tokens, err := lexer.NewLexer().Parse(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	p := parser.NewParser()
	statements, errs := p.Parse(tokens)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		return
	}

	environment := parser.NewEnvironment(nil)
	for _, stmt := range statements {
		err := stmt.Execute(environment)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
	}
}
