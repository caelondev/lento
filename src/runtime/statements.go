package runtime

import "github.com/caelondev/lento/src/ast"

func (i *Interpreter) evaluateBlockStatement(block *ast.BlockStatement) RuntimeValue {
	var lastEvaluated RuntimeValue = NIL()

	for _, statement := range block.Body {
		lastEvaluated = i.EvaluateStatement(statement)
	}

	return lastEvaluated
}
