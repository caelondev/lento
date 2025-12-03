package runtime

import (
	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
)

type Interpreter struct {
	errorHandler *errorhandler.ErrorHandler
	globalEnv    Environment
	line         uint
}

func NewInterpreter(errorHandler *errorhandler.ErrorHandler, env Environment) *Interpreter {
	return &Interpreter{
		errorHandler: errorHandler,
		globalEnv:    env,
		line:         1,
	}
}

func (i *Interpreter) Evaluate(program *ast.BlockStatement) RuntimeValue {
	return i.EvaluateStatement(program, i.globalEnv)
}
