package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/debugg-er/lox/src/lexer"
)

func main() {
	if len(os.Args) > 1 {
		source, err := os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		execute(string(source))

	} else {
		enterPrompt()
	}
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
	scanner := lexer.NewLexer(source)
	tokens := scanner.Parse()
	fmt.Println(tokens)
}
