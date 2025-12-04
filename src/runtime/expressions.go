package runtime

import (
	"fmt"
	"math"

	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
	"github.com/caelondev/lento/src/lexer"
)

func (i *Interpreter) EvaluateExpression(expr ast.Expression, env Environment) RuntimeValue {
	i.line = uint(expr.GetLine())

	switch n := expr.(type) {
	case *ast.NumberExpression:
		return evaluateNumberExpression(n)
	case *ast.StringExpression:
		return evaluateStringExpression(n)
	case *ast.ArrayExpression:
		return i.evaluateArrayExpression(n, env)
	case *ast.BinaryExpression:
		return i.evaluateBinaryExpression(n, env)
	case *ast.UnaryExpression:
		return i.evaluateUnaryExpression(n, env)
	case *ast.SymbolExpression:
		return evaluateSymbolExpression(n, env)
	case *ast.AssignmentExpression:
		return i.evaluateAssignmentExpression(n, env)
	case *ast.CallExpression:
		return i.evaluateCallExpression(n, env)
	case *ast.IndexExpression:
		return i.evaluateIndexExpression(n, env)

	default:
		i.errorHandler.Report(i.line, fmt.Sprintf("Unrecognized AST Expression whilst evaluating: %T\n", expr))
	}

	return NIL()
}

func evaluateNumberExpression(expr *ast.NumberExpression) RuntimeValue {
	return &NumberValue{Value: expr.Value}
}

func evaluateSymbolExpression(expr *ast.SymbolExpression, env Environment) RuntimeValue {
	return env.LookupVariable(expr.Line, expr.Value)
}

func evaluateStringExpression(expr *ast.StringExpression) RuntimeValue {
	value := expr.Value[1 : len(expr.Value)-1]
	return &StringValue{Value: value}
}

func (i *Interpreter) evaluateArrayExpression(expr *ast.ArrayExpression, env Environment) RuntimeValue {
	var elements []RuntimeValue

	for _, element := range expr.Elements {
		elements = append(elements, i.EvaluateExpression(element, env))
	}

	return &ArrayValue{Elements: elements}
}

func (i *Interpreter) evaluateUnaryExpression(expr *ast.UnaryExpression, env Environment) RuntimeValue {
	operand := i.EvaluateExpression(expr.Operand, env)
	operator := expr.Operator.TokenType

	switch operator {
	case lexer.PLUS:
		if num, ok := operand.(*NumberValue); ok {
			return &NumberValue{Value: +num.Value}
		}
		i.errorHandler.Report(i.line, "Unary '+' operator requires a number")
	case lexer.MINUS:
		if num, ok := operand.(*NumberValue); ok {
			return &NumberValue{Value: -num.Value}
		}
		i.errorHandler.Report(i.line, "Unary '-' operator requires a number")
	case lexer.NOT:
		if b, ok := operand.(*BooleanValue); ok {
			return BOOLEAN(!isTruthy(b))
		}
	default:
		i.errorHandler.Report(i.line, fmt.Sprintf("Unrecognized unary operator: %s", expr.Operator.Lexeme))
	}

	return NIL()
}

func (i *Interpreter) evaluateBinaryExpression(expr *ast.BinaryExpression, env Environment) RuntimeValue {
	operatorToken := expr.Operator
	left := i.EvaluateExpression(expr.Left, env)
	right := i.EvaluateExpression(expr.Right, env)

	switch operatorToken.TokenType {
	case lexer.AND:
		return BOOLEAN(isTruthy(left) && isTruthy(right))
	case lexer.OR:
		return BOOLEAN(isTruthy(left) || isTruthy(right))
	default:
		leftNum, leftIsNum := left.(*NumberValue)
		rightNum, rightIsNum := right.(*NumberValue)

		if leftIsNum && rightIsNum {
			return i.evaluateNumericBinaryExpression(leftNum, rightNum, operatorToken)
		}

		leftStr, leftIsStr := left.(*StringValue)
		rightStr, rightIsStr := right.(*StringValue)

		if leftIsStr && rightIsStr {
			return i.evaluateStringBinaryExpression(leftStr, rightStr, operatorToken)
		}

		i.errorHandler.Report(i.line, fmt.Sprintf("Cannot perform '%s' binary operator with unsupported type (%s to %s)", operatorToken.Lexeme, left.Type(), right.Type()))
	}
	return NIL()
}

func (i *Interpreter) evaluateStringBinaryExpression(left *StringValue, right *StringValue, operator *lexer.Token) RuntimeValue {
	lhs := left.Value
	rhs := right.Value

	switch operator.TokenType {
	case lexer.PLUS:
		return &StringValue{Value: lhs + rhs}
	case lexer.EQUAL:
		return BOOLEAN(lhs == rhs)
	case lexer.NOT_EQUAL:
		return BOOLEAN(lhs != rhs)
	default:
		i.errorHandler.Report(i.line, fmt.Sprintf("Unsupported string binary operator: '%s'", operator.Lexeme))
	}

	return NIL()
}

func (i *Interpreter) evaluateNumericBinaryExpression(left *NumberValue, right *NumberValue, operator *lexer.Token) RuntimeValue {
	result := 0.0
	lhs := left.Value
	rhs := right.Value

	switch operator.TokenType {
	case lexer.PLUS:
		result = lhs + rhs
	case lexer.MINUS:
		result = lhs - rhs
	case lexer.STAR:
		result = lhs * rhs
	case lexer.SLASH:
		if rhs == 0 {
			i.errorHandler.Report(i.line, "Division by zero")
		}
		result = lhs / rhs
	case lexer.MODULO:
		if rhs == 0 {
			i.errorHandler.Report(i.line, "Modulo by zero")
		}
		result = math.Mod(lhs, rhs)
	case lexer.LESS:
		return BOOLEAN(lhs < rhs)
	case lexer.LESS_EQUAL:
		return BOOLEAN(lhs <= rhs)
	case lexer.GREATER:
		return BOOLEAN(lhs > rhs)
	case lexer.GREATER_EQUAL:
		return BOOLEAN(lhs >= rhs)
	case lexer.EQUAL:
		return BOOLEAN(lhs == rhs)
	case lexer.NOT_EQUAL:
		return BOOLEAN(lhs != rhs)
	default:
		i.errorHandler.Report(i.line, fmt.Sprintf("Unsupported numeric binary operator: '%s'", operator.Lexeme))
	}

	return &NumberValue{Value: result}
}

func (i *Interpreter) evaluateAssignmentExpression(expr *ast.AssignmentExpression, env Environment) RuntimeValue {
	value := i.EvaluateExpression(expr.Value, env)

	switch assignee := expr.Assignee.(type) {
	case *ast.SymbolExpression:
		env.AssignVariable(i.line, assignee.Value, value)
		return value

	case *ast.IndexExpression:
		array := i.EvaluateExpression(assignee.Array, env)
		index := i.EvaluateExpression(assignee.Index, env)

		arrayValue, isArray := array.(*ArrayValue)
		if !isArray {
			i.errorHandler.ReportError(
				"Interpreter-Array",
				fmt.Sprintf("Cannot index non-array type '%s'", array.Type()),
				i.line,
				errorhandler.ArrayIndexError,
			)
			return NIL()
		}

		indexValue, ok := index.(*NumberValue)
		if !ok {
			i.errorHandler.ReportError(
				"Interpreter-Array",
				fmt.Sprintf("Index must be a number, got '%s'", index.Type()),
				i.line,
				errorhandler.ArrayIndexError,
			)
			return NIL()
		}

		idx := int(indexValue.Value)
		if idx < 0 || idx >= len(arrayValue.Elements) {
			i.errorHandler.ReportError(
				"Interpreter-Array",
				fmt.Sprintf("Index %d out of bounds for array of length %d", idx, len(arrayValue.Elements)),
				i.line,
				errorhandler.ArrayIndexError,
			)
			return NIL()
		}

		arrayValue.Elements[idx] = value
		return value

	default:
		i.errorHandler.Report(i.line, "Invalid left-hand assignment")
		return NIL()
	}
}

func (i *Interpreter) evaluateCallExpression(call *ast.CallExpression, env Environment) RuntimeValue {
	caller := i.EvaluateExpression(call.Caller, env)

	// Parse all string arguments ---
	var args []RuntimeValue
	for _, argExpr := range call.Arguments {
		args = append(args, i.EvaluateExpression(argExpr, env))
	}

	// Native function call ---
	if nativeFunc, ok := caller.(*NativeFunctionValue); ok {
		return nativeFunc.Call(args, env, i)
	}

	// User defined function ---
	if function, ok := caller.(*FunctionValue); ok {
		if len(args) != len(function.Parameters) {
			i.errorHandler.ReportError(
				"Interpreter-Function",
				fmt.Sprintf("Function '%s' expects %d argument(s) but got %d instead", function.Name, len(function.Parameters), len(args)),
				i.line,
				errorhandler.InvalidArgumentError,
			)
			return NIL()
		}

		// Create function scope with the captured environment as parent ---
		functionScope := NewEnvironment(function.Environment, i.errorHandler)

		// Bind parameters to arguments in the function scope ---
		for idx, param := range function.Parameters {
			functionScope.DeclareVariable(i.line, param, args[idx], false, false)
		}

		// Execute body with the function scope
		return i.EvaluateStatement(function.Body, functionScope)
	}

	i.errorHandler.ReportError(
		"Interpreter-Function",
		fmt.Sprintf("Cannot call non-function expression type '%s'", caller.Type()),
		i.line,
		errorhandler.NonFunctionExpressionError,
	)
	return NIL()
}

func (i *Interpreter) evaluateIndexExpression(expr *ast.IndexExpression, env Environment) RuntimeValue {
	array := i.EvaluateExpression(expr.Array, env)
	index := i.EvaluateExpression(expr.Index, env)

	arrayValue, isArray := array.(*ArrayValue)
	if !isArray {
		i.errorHandler.ReportError(
			"Interpreter-Array",
			fmt.Sprintf("Cannot index expression as it is not an array (type of '%s')", array.Type()),
			i.line,
			errorhandler.ArrayIndexError,
		)
		return NIL()
	}

	indexValue, ok := index.(*NumberValue)
	if !ok {
		i.errorHandler.ReportError(
			"Interpreter-Array",
			fmt.Sprintf("Invalid indexing value (type of %s)", array.Type()),
			i.line,
			errorhandler.ArrayIndexError,
		)
		return NIL()
	}

	if len(arrayValue.Elements) < int(indexValue.Value) {
		i.errorHandler.ReportError(
			"Interpreter-Array",
			fmt.Sprintf("Index %d out of bounds as the array only got %d elements", int(indexValue.Value), len(arrayValue.Elements)),
			i.line,
			errorhandler.ArrayIndexError,
		)
		return NIL()
	}

	for idx, array := range arrayValue.Elements {
		if idx == int(indexValue.Value) {
			return array
		}
	}

	return NIL()
}
