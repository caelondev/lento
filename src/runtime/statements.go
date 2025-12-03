package runtime

import "github.com/caelondev/lento/src/ast"

func (i *Interpreter) evaluateBlockStatement(block *ast.BlockStatement) RuntimeValue {
	previousEnv := i.GetEnvironment()
	i.environment = i.EnterScope()

	defer func() { i.environment = previousEnv }()

	var lastEvaluated RuntimeValue = NIL()
	for _, statement := range block.Body {
		lastEvaluated = i.EvaluateStatement(statement)
	}

	return lastEvaluated
}

func (i *Interpreter) evaluateVariableDeclarationStatement(decl *ast.VariableDeclarationStatement) RuntimeValue {
	var value RuntimeValue = NIL()

	env := i.GetEnvironment()
	if decl.Value != nil {
		value = i.EvaluateExpression(decl.Value)
	}

	env.DeclareVariable(i.line, decl.Identifier, value, decl.IsConstant, false)
	return value
}

func (i *Interpreter) evaluateIfStatement(stmt *ast.IfStatement) RuntimeValue {
	condition := i.EvaluateExpression(stmt.Condition)

	if isTruthy(condition) {
		i.EvaluateStatement(stmt.Consequent)
	} else {
		i.EvaluateStatement(stmt.Alternate)
	}

	return NIL()
}

func (i *Interpreter) evaluateFunctionDeclaration(stmt *ast.FunctionDeclarationStatement) RuntimeValue {
	env := i.GetEnvironment()

	fn := &FunctionValue{
		Name: stmt.Name,
		Parameters: stmt.Parameters,
		Body: stmt.Body,
		Environment: env,
	}

	env.DeclareFunction(i.line, stmt.Name, fn, false)
	return fn
}
