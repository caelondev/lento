package src

import (
	"fmt"
	"os"

	errorhandler "github.com/caelondev/lento/src/error-handler"
	"github.com/caelondev/lento/src/lexer"
)

var ErrorHandler = errorhandler.New()

func Lento() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lento [filepath]")
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1])
	}
	
}

func runFile(filepath string) {
	bytes, error := os.ReadFile(filepath)
	if error != nil {
		fmt.Printf("An error occurred whilst trying to read %s:\n%s\n", filepath, error.Error())
	}

	run(string(bytes))
}

func run(sourceCode string) {
	lexer := lexer.NewLexer(sourceCode, ErrorHandler)
	tokens := lexer.Tokenize()

	for _, t := range tokens {
		t.String()
	}
}
