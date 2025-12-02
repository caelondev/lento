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

func (i *Interpreter) evaluateStringExpression(expr *ast.StringExpression) RuntimeValue {
	value := expr.Value[1 : len(expr.Value) - 1] // Trim "
	return &StringValue{Value: value}
}

func (i *Interpreter) evaluateUnaryExpression(expr *ast.UnaryExpression) RuntimeValue {
	operand := i.EvaluateExpression(expr.Operand)
	operator := expr.Operator.TokenType

	switch operator {
	case lexer.PLUS:
		if num, ok := operand.(*NumberValue); ok {
			return &NumberValue{Value: +num.Value}
		}
		i.errorHandler.Report(int(i.line), "Unary '+' operator requires a number")
	case lexer.MINUS:
		if num, ok := operand.(*NumberValue); ok {
			return &NumberValue{Value: -num.Value}
		}
		i.errorHandler.Report(int(i.line), "Unary '-' operator requires a number")
	case lexer.NOT:
		if b, ok := operand.(*BooleanValue); ok {
			return BOOLEAN(!isTruthy(b))
		}

	default:
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unrecognized unary operator: %s", expr.Operator.Lexeme))
	}

	return NIL()
}

func (i *Interpreter) evaluateBinaryExpression(expr *ast.BinaryExpression) RuntimeValue {
	operatorToken := expr.Operator
	left := i.EvaluateExpression(expr.Left)
	right := i.EvaluateExpression(expr.Right)

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

		i.errorHandler.Report(int(i.line), fmt.Sprintf("Cannot perform '%s' binary operator with unsupported type (%s to %s)", operatorToken.Lexeme, left.Type(), right.Type()))
	}
	return NIL()
}

func (i *Interpreter) evaluateStringBinaryExpression(left *StringValue, right *StringValue, operator *lexer.Token) RuntimeValue {
	var value string

	lhs := left.Value
	rhs := right.Value

	switch operator.TokenType {
	case lexer.PLUS:
		value = lhs + rhs
	case lexer.EQUAL:
		return BOOLEAN(lhs == rhs)
	case lexer.NOT_EQUAL:
		return BOOLEAN(lhs != rhs)

	default:
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unsupported string binary operator: '%s'", operator.Lexeme))
	}

	return &StringValue{Value: value}
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
		i.errorHandler.Report(int(i.line), fmt.Sprintf("Unsupported numeric binary operator: '%s'", operator.Lexeme))
	}

	return &NumberValue{Value: result}
}

