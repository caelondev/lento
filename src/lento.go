package src

import (
	"bufio"
	"fmt"
	"os"

	errorhandler "github.com/caelondev/lento/src/error-handler"
	"github.com/caelondev/lento/src/lexer"
	"github.com/caelondev/lento/src/parser"
	"github.com/caelondev/lento/src/runtime"
	"github.com/sanity-io/litter"
)

var ErrorHandler = errorhandler.New()

func Lento() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lento [filepath]")
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runRepl()
	}

}

func runFile(filepath string) {
	bytes, error := os.ReadFile(filepath)
	if error != nil {
		fmt.Printf("An error occurred whilst trying to read %s:\n%s\n", filepath, error.Error())
	}

	run(string(bytes))

	if ErrorHandler.HadError {
		os.Exit(ErrorHandler.ErrorCode)
	}
}

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		run(line)
		ErrorHandler.HadError = false
	}
}

func run(sourceCode string) {
	lexer := lexer.NewLexer(sourceCode, ErrorHandler)
	interpreter := runtime.NewInterpreter(ErrorHandler)

	tokens := lexer.Tokenize()
	if ErrorHandler.HadError {
		return
	}

	ast := parser.ProduceAST(tokens, ErrorHandler)
	if ErrorHandler.HadError {
		return
	}

	for _, token := range tokens {
		token.String()
	}

	litter.Dump(ast)
	
	var result runtime.RuntimeValue
	for _, statement := range ast.Body {
		result = interpreter.EvaluateStatement(statement)
		if ErrorHandler.HadError {
			return
		}
	}

	litter.Dump(result)
}
