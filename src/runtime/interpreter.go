package runtime

import (
	"fmt"

	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
)

type Interpreter struct {
	errorHandler *errorhandler.ErrorHandler
	line         uint
}

func NewInterpreter(errorHandler *errorhandler.ErrorHandler) *Interpreter {
	return &Interpreter{
		errorHandler: errorHandler,
		line:         1,
	}
}

func (i *Interpreter) EvaluateStatement(node ast.Statement) RuntimeValue {
	switch n := node.(type) {
	case *ast.BlockStatement:
		return i.evaluateBlockStatement(n)
	case *ast.ExpressionStatement:
		return i.EvaluateExpression(n.Expression)

	default:
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unrecognized AST Statement whilst evaluating: %T\n", node))

	}

	return nil
}

func (i *Interpreter) EvaluateExpression(node ast.Expression) RuntimeValue {
	switch n := node.(type) {
	case *ast.NumberExpression:
		return i.evaluateNumberExpression(n)
	case *ast.BinaryExpression:
		return i.evaluateBinaryExpression(n)

	default:
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unrecognized AST Expression whilst evaluating: %T\n", node))
	}

	return nil
}
