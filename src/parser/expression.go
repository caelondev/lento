package parser

import (
	"fmt"
	"strconv"

	"github.com/caelondev/lento/src/ast"
	"github.com/caelondev/lento/src/lexer"
)

func parseExpression(p *parser, bp BindingPower) ast.Expression {
	// Parse NUD ---
	
	tokenType := p.currentTokenType()
	nudFunction, exists := nudLU[tokenType]

	if p.isEOF() {
		p.errorHandler.ReportError(
			"Parser-NUD",
			"Unexpected end of file expression",
			int(p.line),
			65,
		)
	}

	if !exists {
		p.errorHandler.ReportError(
			"Parser-NUD",
			fmt.Sprintf("Unrecognized token found whilst parsing: '%s'", lexer.TokenTypeString[tokenType]),
			int(p.line),
			65,
		)
	}

	left := nudFunction(p)

	for !p.isEOF() && bindingPowerLU[p.currentTokenType()] > bp {
		tokenType = p.currentTokenType()
		ledFunction, exists := ledLU[tokenType]

		if !exists {
			p.errorHandler.ReportError(
				"Parser-LED",
				fmt.Sprintf("Unrecognized token found in the middle of an expression: '%s'", lexer.TokenTypeString[tokenType]),
				int(p.line),
				65,
			)
		}

		left = ledFunction(p, left, bindingPowerLU[p.currentTokenType()])
	}

	return left
}

func parsePrimaryExpression(p *parser) ast.Expression {
	switch p.currentTokenType() {
	case lexer.NUMBER:
		number, err := strconv.ParseFloat(p.advance().Lexeme, 64)
		if err != nil {
			p.errorHandler.Report(int(p.line), fmt.Sprintf("An Error occured whilst trying to parse a number:\n%s", err))
		}

		return &ast.NumberExpression{
			Value: number,
		}
	case lexer.STRING:
		return &ast.StringExpression{
			Value: p.advance().Lexeme,
		}
	case lexer.IDENTIFIER:
		return &ast.SymbolExpression{
			Value: p.advance().Lexeme,
		}
	case lexer.LEFT_PARENTHESIS:
		p.advance() // Eat '(' ---
		value := parseExpression(p, DEFAULT_BP)
		p.expect(lexer.RIGHT_PARENTHESIS)
		return value

	default:
		p.errorHandler.Report(
			int(p.line),
			fmt.Sprintf("Unrecognized primary token: %s", lexer.TokenTypeString[p.currentTokenType()]),
		)
	}

	return nil
}

func parseBinaryExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance() // Eat 
	right := parseExpression(p, DEFAULT_BP)

	return &ast.BinaryExpression{
		Left: left,
		Right: right,
		Operator: operatorToken,
	}
}
