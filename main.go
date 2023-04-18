package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/debugg-er/lox/src/interpreter"
	"github.com/debugg-er/lox/src/lexer"
	"github.com/debugg-er/lox/src/parser"
)

func main() {
	start := time.Now()
	if len(os.Args) > 1 {
		ExecFile()
	} else {
		EnterPrompt()
	}
	fmt.Printf("Program exited after %dms...\n", time.Since(start).Milliseconds())
}

func ExecFile() {
	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "File not found")
		os.Exit(1)
	}
	execute(string(source))
}

func EnterPrompt() {
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
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return
	}

	interpreter := interpreter.NewInterpreter()
	if err := interpreter.Run(statements); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}
