package parser

import (
	"fmt"

	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
	"github.com/caelondev/lento/src/lexer"
)

func parseStatement(p *parser) ast.Statement {
	if p.isEOF() {
		return nil
	}

	statementFunction, exists := statementLU[p.currentTokenType()]

	if exists {
		stmt := statementFunction(p)
		return stmt
	}

	expression := parseExpression(p, DEFAULT_BP)

	if expression == nil {
		p.synchronize()
		return nil // Error already reported
	}

	if p.currentTokenType() != lexer.SEMICOLON {
		p.errorHandler.ReportError(
			"Parser",
			fmt.Sprintf("Expected SEMICOLON after expression, got %s instead", lexer.TokenTypeString[p.currentTokenType()]),
			p.line,
			errorhandler.UnexpectedTokenError,
		)
		p.synchronize()
		return nil
	}

	p.expect(lexer.SEMICOLON)

	return &ast.ExpressionStatement{
		Expression: expression,
		Line:       p.line,
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
		Value:      value,
		Line:       p.line,
	}
}

func parseBlockStatement(p *parser) ast.Statement {
	statements := make([]ast.Statement, 0)

	for !p.isEOF() && p.currentTokenType() != lexer.RIGHT_BRACE {
		stmt := parseStatement(p)
		if stmt == nil {
			break
		}
		statements = append(statements, stmt)
	}

	p.expect(lexer.RIGHT_BRACE)

	return &ast.BlockStatement{
		Body: statements,
		Line: p.line,
	}
}

func parseIfStatement(p *parser) ast.Statement {

	// SYNTAX ---
	//
	// if (condition) { ... }
	// if (condition) ...
	// if condition { ... }
	//
	// else { ... }
	// else ...
	//

	var condition ast.Expression
	var consequent ast.Statement // If block ---
	var alternate ast.Statement  // else block

	p.expect(lexer.IF)

	p.expect(lexer.LEFT_PARENTHESIS)
	condition = parseExpression(p, DEFAULT_BP)
	p.expect(lexer.RIGHT_PARENTHESIS)

	if p.currentTokenType() == lexer.LEFT_BRACE {
		p.advance() // Eat brace ---
		consequent = parseBlockStatement(p)
	} else {
		consequent = parseStatement(p)
	}

	// Check for 'else' token ---
	if p.currentTokenType() == lexer.ELSE {
		p.advance() // Eat 'else' ---

		if p.currentTokenType() == lexer.LEFT_BRACE {
			p.advance() // Eat '}' ---
			alternate = parseBlockStatement(p)
		} else {
			alternate = parseStatement(p)
		}
	}

	return &ast.IfStatement{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
		Line:       p.line,
	}
}

func parseFunctionDeclaration(p *parser) ast.Statement {
	// SYNTAX ---
	// fn identifier(params) { ... }
	// fn identifier(params) ...
	//

	var identifier string
	var parameters []string
	var body ast.Statement

	p.expect(lexer.FUNCTION)

	identifier = p.expect(lexer.IDENTIFIER).Lexeme

	p.expect(lexer.LEFT_PARENTHESIS)

	if p.currentTokenType() != lexer.RIGHT_PARENTHESIS {
		param := p.expect(lexer.IDENTIFIER).Lexeme
		parameters = append(parameters, param)

		// Parse remaining parameters (comma-separated)
		for p.currentTokenType() == lexer.COMMA {
			p.advance() // eat comma
			param := p.expect(lexer.IDENTIFIER).Lexeme
			parameters = append(parameters, param)
		}
	}

	p.expect(lexer.RIGHT_PARENTHESIS)

	if p.currentTokenType() == lexer.LEFT_BRACE {
		p.advance()
		body = parseBlockStatement(p)
	} else {
		body = parseStatement(p)
	}

	return &ast.FunctionDeclarationStatement{
		Name:       identifier,
		Parameters: parameters,
		Body:       body,
		Line:       p.line,
	}
}

func parseWhileStatement(p *parser) ast.Statement {
	// SYNTAX ---
	// while condition { ... }
	// while (condition) { ... }
	// while (condition) ...;
	//

	var condition ast.Expression
	var body ast.Statement

	p.expect(lexer.WHILE)

	// Optional parentheses around condition
	if p.currentTokenType() == lexer.LEFT_PARENTHESIS {
		p.advance() // Eat '('
		condition = parseExpression(p, DEFAULT_BP)
		p.expect(lexer.RIGHT_PARENTHESIS)
	} else {
		condition = parseExpression(p, DEFAULT_BP)
	}

	// Parse the loop body
	if p.currentTokenType() == lexer.LEFT_BRACE {
		p.advance() // Eat '{'
		body = parseBlockStatement(p)
	} else {
		// Single statement loop body (no braces)
		body = parseStatement(p)
	}

	return &ast.WhileLoopStatement{
		Condition: condition,
		Body:      body,
		Line:      p.line,
	}
}

func parseForStatement(p *parser) ast.Statement {
	// SYNTAX ---
	//
	// for (var x = 0; x<10; x++) { ... }
	// for (var x = 0; x<10; x++) ...
	//
	// INFINITE LOOP ---
	// for (;;) { ... }
	// for (;;) ...
	//

	var init ast.Statement
	var condition ast.Expression
	var increment ast.Expression
	var body ast.Statement

	p.expect(lexer.FOR)

	p.expect(lexer.LEFT_PARENTHESIS)

	if p.currentTokenType() != lexer.VARIABLE {
		p.errorHandler.ReportError(
			"Parser-For",
			fmt.Sprintf("Expected a non-constant variable declaration, got %s instead", lexer.TokenTypeString[p.currentTokenType()]),
			p.line,
			errorhandler.VariableDeclarationError,
		)
		return nil
	}

	init = parseVariableDeclaration(p) // Already requires ';' ---

	condition = parseExpression(p, DEFAULT_BP)
	p.expect(lexer.SEMICOLON)

	increment = parseExpression(p, DEFAULT_BP)
	p.expect(lexer.RIGHT_PARENTHESIS)

	if p.currentTokenType() == lexer.LEFT_BRACE {
		p.advance()
		body = parseBlockStatement(p)
	} else {
		body = parseStatement(p)
	}

	return &ast.ForStatement{
		Init:      init,
		Condition: condition,
		Increment: increment,
		Body:      body,
		Line:      p.line,
	}
}
