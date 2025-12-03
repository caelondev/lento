package runtime

import (
	"fmt"

	"github.com/caelondev/lento/src/ast"
)

func (i *Interpreter) EvaluateStatement(stmt ast.Statement, env Environment) RuntimeValue {
	i.line = uint(stmt.GetLine())

	switch n := stmt.(type) {
	case *ast.BlockStatement:
		return i.evaluateBlockStatement(n, env)
	case *ast.ExpressionStatement:
		return i.EvaluateExpression(n.Expression, env)
	case *ast.VariableDeclarationStatement:
		return evaluateVariableDeclarationStatement(n, env, i)
	case *ast.IfStatement:
		return i.evaluateIfStatement(n, env)
	case *ast.FunctionDeclarationStatement:
		return evaluateFunctionDeclaration(n, env)
	default:
		i.errorHandler.Report(i.line, fmt.Sprintf("Unrecognized AST Statement whilst evaluating: %T\n", stmt))
	}

	return NIL()
}

func (i *Interpreter) evaluateBlockStatement(block *ast.BlockStatement, env Environment) RuntimeValue {
	// Create new scope
	blockScope := NewEnvironment(env, i.errorHandler)

	var lastEvaluated RuntimeValue = NIL()
	for _, statement := range block.Body {
		lastEvaluated = i.EvaluateStatement(statement, blockScope)
	}

	return lastEvaluated
}

func evaluateVariableDeclarationStatement(decl *ast.VariableDeclarationStatement, env Environment, i *Interpreter) RuntimeValue {
	var value RuntimeValue = NIL()

	if decl.Value != nil {
		value = i.EvaluateExpression(decl.Value, env)
	}

	env.DeclareVariable(i.line, decl.Identifier, value, decl.IsConstant, false)
	return value
}

func (i *Interpreter) evaluateIfStatement(stmt *ast.IfStatement, env Environment) RuntimeValue {
	condition := i.EvaluateExpression(stmt.Condition, env)

	if isTruthy(condition) {
		i.EvaluateStatement(stmt.Consequent, env)
	} else {
		i.EvaluateStatement(stmt.Alternate, env)
	}

	return NIL()
}

func evaluateFunctionDeclaration(stmt *ast.FunctionDeclarationStatement, env Environment) RuntimeValue {
	// Capture the current environment (closure)
	fn := &FunctionValue{
		Name:        stmt.Name,
		Parameters:  stmt.Parameters,
		Body:        stmt.Body,
		Environment: env,
	}

	env.DeclareVariable(stmt.Line, stmt.Name, fn, true, false)
	return fn
}
