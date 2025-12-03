package runtime

import (
	"fmt"

	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
)

type Interpreter struct {
	errorHandler *errorhandler.ErrorHandler
	environment  Environment
	line         uint
}

func NewInterpreter(errorHandler *errorhandler.ErrorHandler, env Environment) *Interpreter {
	return &Interpreter{
		errorHandler: errorHandler,
		environment:  env,
		line:         1,
	}
}

func (i *Interpreter) EvaluateStatement(node ast.Statement) RuntimeValue {
	// Update line before processing
	i.line = uint(node.GetLine())

	switch n := node.(type) {
	case *ast.BlockStatement:
		return i.evaluateBlockStatement(n)
	case *ast.ExpressionStatement:
		return i.EvaluateExpression(n.Expression)
	case *ast.VariableDeclarationStatement:
		return i.evaluateVariableDeclarationStatement(n)
	case *ast.IfStatement:
		return i.evaluateIfStatement(n)
	case *ast.FunctionDeclarationStatement:
		return i.evaluateFunctionDeclaration(n)

	default:
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unrecognized AST Statement whilst evaluating: %T\n", node))
	}

	return NIL()
}

func (i *Interpreter) EvaluateExpression(node ast.Expression) RuntimeValue {
	// Update line before processing
	i.line = uint(node.GetLine())

	switch n := node.(type) {
	case *ast.NumberExpression:
		return i.evaluateNumberExpression(n)
	case *ast.StringExpression:
		return i.evaluateStringExpression(n)
	case *ast.BinaryExpression:
		return i.evaluateBinaryExpression(n)
	case *ast.UnaryExpression:
		return i.evaluateUnaryExpression(n)
	case *ast.SymbolExpression:
		return i.evaluateSymbolExpression(n)
	case *ast.AssignmentExpression:
		return i.evaluateAssignmentExpression(n)

	default:
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unrecognized AST Expression whilst evaluating: %T\n", node))
	}

	return NIL()
}

func (i *Interpreter) GetEnvironment() Environment {
	return i.environment
}

func (i *Interpreter) EnterScope() Environment {
	newEnv := NewEnvironment(i.GetEnvironment(), i.errorHandler)
	return newEnv
}
