package parser

import (
	"github.com/caelondev/lento/src/ast"
	"github.com/caelondev/lento/src/lexer"
)

func parseStatement(p *parser) ast.Statement {
	if p.isEOF() {
		return nil
	}

	statementFunction, exists := statementLU[p.currentTokenType()]

	if exists {
		return statementFunction(p)
	}

	expression := parseExpression(p, DEFAULT_BP)

	p.expect(lexer.SEMICOLON)

	return &ast.ExpressionStatement{
		Expression: expression,
		Line: p.line,
	}
}

func parseVariableDeclaration(p *parser) ast.Statement {
	//
	//  var <identifier> = [value]; ---
	//  var <identifier>;
	//	const <identifier> = <value>; ---
	//
	
	isConstant := p.advance().TokenType == lexer.CONSTANT
	identifier := p.expect(lexer.IDENTIFIER).Lexeme
	var value ast.Expression

	if p.currentTokenType() != lexer.SEMICOLON {
		p.expect(lexer.ASSIGNMENT, lexer.SEMICOLON) // NOTE: SEMICOLON IS NOT NEEDED HERE, I JUST ADDED IT FOR ERROR MESSAGE --
		value = parseExpression(p, DEFAULT_BP)
	}

	p.expect(lexer.SEMICOLON)

	return &ast.VariableDeclarationStatement{
		IsConstant: isConstant,
		Identifier: identifier,
		Value: value,
		Line: p.line,
	}
}
