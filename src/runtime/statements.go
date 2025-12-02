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
	env := i.GetEnvironment()
	value := i.EvaluateExpression(decl.Value)
	env.DeclareVariable(i.line, decl.Identifier, value, decl.IsConstant, false)
	return value
}
