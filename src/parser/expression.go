package parser

import (
	"fmt"
	"strconv"

	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
	"github.com/caelondev/lento/src/lexer"
)

func parseExpression(p *parser, bp BindingPower) ast.Expression {
	tokenType := p.currentTokenType()
	nudFunction, exists := nudLU[tokenType]

	if p.isEOF() {
		p.errorHandler.ReportError(
			"Parser",
			"Unexpected end of input while parsing expression",
			p.line,
			errorhandler.UnexpectedTokenError,
		)
		return nil
	}

	if !exists {
		p.errorHandler.ReportError(
			"Parser",
			fmt.Sprintf("Unexpected token '%s' - expected an expression", p.currentToken().Lexeme),
			p.line,
			errorhandler.UnknownTokenError,
		)
		p.synchronize() // Skip to next safe point ---
		return nil
	}

	left := nudFunction(p)
	if left == nil {
		return nil // Propagate error ---
	}

	for !p.isEOF() && bindingPowerLU[p.currentTokenType()] > bp {
		tokenType = p.currentTokenType()
		operatorBP := bindingPowerLU[tokenType]
		ledFunction, exists := ledLU[tokenType]

		if !exists {
			p.errorHandler.ReportError(
				"Parser",
				fmt.Sprintf("Unexpected token '%s' in expression", p.currentToken().Lexeme),
				p.line,
				errorhandler.UnexpectedTokenError,
			)
			p.synchronize()
			return left // Return what we have so far ---
		}

		left = ledFunction(p, left, operatorBP)
		if left == nil {
			return nil
		}
	}

	return left
}

func parsePrimaryExpression(p *parser) ast.Expression {
	switch p.currentTokenType() {
	case lexer.NUMBER:
		number, err := strconv.ParseFloat(p.advance().Lexeme, 64)
		if err != nil {
			p.errorHandler.Report(p.line, fmt.Sprintf("An Error occured whilst trying to parse a number:\n%s", err))
		}

		return &ast.NumberExpression{
			Value: number,
			Line:  p.line,
		}
	case lexer.STRING:
		return &ast.StringExpression{
			Value: p.advance().Lexeme,
			Line:  p.line,
		}
	case lexer.IDENTIFIER:
		return &ast.SymbolExpression{
			Value: p.advance().Lexeme,
			Line:  p.line,
		}
	case lexer.LEFT_PARENTHESIS:
		p.advance() // Eat '(' ---
		value := parseExpression(p, DEFAULT_BP)
		p.expect(lexer.RIGHT_PARENTHESIS)
		return value

	default:
		p.errorHandler.Report(
			p.line,
			fmt.Sprintf("Unrecognized primary token: %s", lexer.TokenTypeString[p.currentTokenType()]),
		)
	}

	return nil
}

func parseBinaryExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance() // Eat operator
	right := parseExpression(p, bp)

	return &ast.BinaryExpression{
		Left:     left,
		Right:    right,
		Operator: operatorToken,
		Line:     p.line,
	}
}

func parseUnaryExpression(p *parser) ast.Expression {
	operatorToken := p.advance()
	value := parseExpression(p, UNARY)
	return &ast.UnaryExpression{
		Operator: operatorToken,
		Operand:  value,
		Line:     p.line,
	}
}

func parseAssignmentExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.advance()
	value := parseExpression(p, ASSIGNMENT-1)

	return &ast.AssignmentExpression{
		Assignee: left,
		Value:    value,
	}
}

func parseCallExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.advance() // Eat '(' ---

	arguments := make([]ast.Expression, 0)

	// Parse arguments (comma-separated expressions)
	if p.currentTokenType() != lexer.RIGHT_PARENTHESIS {
		// Parse first argument
		arg := parseExpression(p, DEFAULT_BP)
		if arg != nil {
			arguments = append(arguments, arg)
		}

		// Parse remaining arguments
		for p.currentTokenType() == lexer.COMMA {
			p.advance() // eat comma
			arg := parseExpression(p, DEFAULT_BP)
			if arg != nil {
				arguments = append(arguments, arg)
			}
		}
	}

	p.expect(lexer.RIGHT_PARENTHESIS)

	return &ast.CallExpression{
		Caller:    left,
		Arguments: arguments,
		Line:      p.line,
	}
}

func parseArrayExpression(p *parser) ast.Expression {
	var elements []ast.Expression

	p.advance() // Eat LEFT_BRACKET token ---

	if p.currentTokenType() != lexer.RIGHT_BRACKET {
		element := parseExpression(p, DEFAULT_BP)
		if element != nil {
			elements = append(elements, element)
		}

		for p.currentTokenType() == lexer.COMMA {
			p.advance() // Eat COMMA
			element := parseExpression(p, DEFAULT_BP)
			if element != nil {
				elements = append(elements, element)
			}
		}
	}

	p.expect(lexer.RIGHT_BRACKET)

	return &ast.ArrayExpression{
		Elements: elements,
		Line:     p.line,
	}
}

func parseIndexExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.advance() // Eat LEFT_PARENTHESIS ---
	index := parseExpression(p, DEFAULT_BP)

	p.expect(lexer.RIGHT_BRACKET)

	return &ast.IndexExpression{
		Array: left,
		Index: index,
		Line:  p.line,
	}
}
