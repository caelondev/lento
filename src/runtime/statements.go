package runtime

import (
	"fmt"

	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
)

func (i *Interpreter) EvaluateStatement(stmt ast.Statement, env Environment) RuntimeValue {
	if i.errorHandler.HadError {
		return NIL()
	}

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
	case *ast.WhileLoopStatement:
		return i.evaluateWhileLoopStatement(n, env)
	case *ast.ForStatement:
		return i.evaluateForStatement(n, env)
	case *ast.ReturnStatement:
		return i.evaluateReturnStatement(n, env)
	case *ast.BreakStatement:
		return i.evaluateBreakStatement(n, env)
	case *ast.ContinueStatement:
		return i.evaluateContinueStatement(n, env)


	default:
		i.errorHandler.Report(i.line, fmt.Sprintf("Unrecognized AST Statement whilst evaluating: %T\n", stmt))
	}

	return NIL()
}

func (i *Interpreter) evaluateBlockStatement(block *ast.BlockStatement, env Environment) RuntimeValue {
	// Create new scope
	blockScope := NewEnvironment(env, i.errorHandler)

	// THIS CODE IS FOR DEBUGGING PURPOSES,UNCOMMENT THIS TO PRINT EVERY ENV CREATED IN REAL TIME ---
	//
	// if envStruct, ok := blockScope.(*EnvironmentStruct); ok {
	// 	fmt.Println("\n=== ENTERING BLOCK SCOPE ===")
	// 	envStruct.Debug(0)
	// 	fmt.Println("============================")
	// }

	var lastEvaluated RuntimeValue = NIL()
	for _, statement := range block.Body {
		lastEvaluated = i.EvaluateStatement(statement, blockScope)

		if _, ok := lastEvaluated.(*ControlFlowValue); ok {
			return lastEvaluated
		}
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
        if stmt.Consequent != nil {
            return i.EvaluateStatement(stmt.Consequent,env)
        }
    } else if stmt.Alternate != nil {
        return i.EvaluateStatement(stmt.Alternate, env)
    }
    return nil
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

func (i *Interpreter) evaluateWhileLoopStatement(stmt *ast.WhileLoopStatement, env Environment) RuntimeValue {
	wasInLoop := i.isInLoop
	i.isInLoop = true

	defer func() { i.isInLoop = wasInLoop }()

	condition := i.EvaluateExpression(stmt.Condition, env)

	for isTruthy(condition) {
		result := i.EvaluateStatement(stmt.Body, env)

		if control, ok := result.(*ControlFlowValue); ok {
			switch control.FlowType {
			case FLOW_BREAK:
				return NIL()
			case FLOW_CONTINUE:
				// Continue
			case FLOW_RETURN:
				return control  // Propagate return up
			}
		}

		condition = i.EvaluateExpression(stmt.Condition, env)
	}

	return NIL()
}

func (i *Interpreter) evaluateForStatement(stmt *ast.ForStatement, env Environment) RuntimeValue {
	forScope := NewEnvironment(env, i.errorHandler)

	wasInLoop := i.isInLoop
	i.isInLoop = true

	defer func() { i.isInLoop = wasInLoop }()

	i.EvaluateStatement(stmt.Init, forScope) // Initialize initializer variable

	for {
		if stmt.Condition != nil {
			condition := i.EvaluateExpression(stmt.Condition, forScope)
			if !isTruthy(condition) {
				break
			}
		}

		result := i.EvaluateStatement(stmt.Body, forScope)

		if control, ok := result.(*ControlFlowValue); ok {
			switch control.FlowType {
			case FLOW_BREAK:
				return NIL()
			case FLOW_CONTINUE:
				// Continue
			case FLOW_RETURN:
				return control  // Propagate return up
			}
		}

		i.EvaluateExpression(stmt.Increment, forScope) // Increment initializer
	}

	return NIL()
}

func (i *Interpreter) evaluateReturnStatement(stmt *ast.ReturnStatement, env Environment) RuntimeValue {
	var value RuntimeValue = NIL()

	if !i.isInFunction {
		i.errorHandler.ReportError(
			"Interpreter-Return",
			"Illegal return statement outside of a function body",
			i.line,
			errorhandler.IllegalStatementError,
		)
	}

	if stmt.Value != nil {
		value = i.EvaluateExpression(stmt.Value, env)
	}

	return RETURN(value)
}

func (i *Interpreter) evaluateBreakStatement(stmt *ast.BreakStatement, env Environment) RuntimeValue {
	if !i.isInLoop {
		i.errorHandler.ReportError(
			"Interpreter-Break",
			"Illegal break statement outside of a loop body",
			i.line,
			errorhandler.IllegalStatementError,
		)
	}

	return BREAK()
}

func (i *Interpreter) evaluateContinueStatement(stmt *ast.ContinueStatement, env Environment) RuntimeValue {
	if !i.isInLoop {
		i.errorHandler.ReportError(
			"Interpreter-Continue",
			"Illegal continue statement outside of a loop body",
			i.line,
			errorhandler.IllegalStatementError,
		)
	}

	return CONTINUE()
}
