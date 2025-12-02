package runtime

import (
	"fmt"
	"math"

	"github.com/caelondev/lento/src/ast"
	"github.com/caelondev/lento/src/lexer"
)

func (i *Interpreter) evaluateNumberExpression(expr *ast.NumberExpression) RuntimeValue {
	return &NumberValue{Value: expr.Value}
}

func (i *Interpreter) evaluateBinaryExpression(expr *ast.BinaryExpression) RuntimeValue {
	left := i.EvaluateExpression(expr.Left)
	right := i.EvaluateExpression(expr.Right)

	leftNum, leftIsNum := left.(*NumberValue)
	rightNum, rightIsNum := right.(*NumberValue)

	if leftIsNum && rightIsNum {
		return i.evaluateNumericBinaryExpression(leftNum, rightNum, expr.Operator)
	}

	i.errorHandler.Report(int(i.line), fmt.Sprintf("Cannot perform %s binary operator with unsupported type (%s to %s)", expr.Operator.Lexeme, left.Type(), right.Type()))
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
			i.errorHandler.Report(int(i.line), "Division by zero")
		}
		result = lhs / rhs
	case lexer.MODULO:
		if rhs == 0 {
			i.errorHandler.Report(int(i.line), "Modulo by zero")
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
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unsupported binary operator: '%s'", operator.Lexeme))
	}

	return &NumberValue{Value: result}
}

